package main

/*

Data is stored using boltdb. We use two buckets, one for storing
metadata and the other for storing chunks. We chunk data to reduce
memory requirements on the server. For any given key β, if we store
the value in ω chunks, the value of β in the meta bucket will be ω,
and we will store chunk 0 as 0~β, 1 as, 1~β, up through (ω-1)~β in
the chunk bucket.

*/

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/grpclog/glogger"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/uluyol/mudahkv/pb"
)

var (
	keyBucketName  = []byte("key")
	metaBucketName = []byte("meta")

	dbPath  = flag.String("f", "", "path to database file")
	portNum = flag.Int("p", 6070, "port to listen on")
)

type server struct {
	db *bolt.DB
}

func newServer(dbPath string) (*server, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	if _, err := tx.CreateBucketIfNotExists(keyBucketName); err != nil {
		tx.Rollback()
		return nil, err
	}
	if _, err := tx.CreateBucketIfNotExists(metaBucketName); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &server{db}, nil
}

type chunkSender interface {
	Send(*pb.KVChunk) error
}

func (s *server) getChunks(tx *bolt.Tx, β string, stream chunkSender) error {
	ωBytes := tx.Bucket(metaBucketName).Get([]byte(β))
	if ωBytes == nil {
		return errors.New("no such file for key")
	}
	ω, err := strconv.Atoi(string(ωBytes))
	if err != nil {
		return fmt.Errorf("failed to decode chunk count, possible corruption: %v", err)
	}
	for i := 0; i < ω; i++ {
		c := tx.Bucket(keyBucketName).Get([]byte(fmt.Sprintf("%d~%s", i, β)))
		if err := stream.Send(&pb.KVChunk{Key: β, Value: c}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) Get(req *pb.Key, stream pb.Mudah_GetServer) error {
	glog.Infof("GET %s", req.Key)
	tx, err := s.db.Begin(false)
	if err != nil {
		return fmt.Errorf("failed to initiate tx: %v", err)
	}
	defer tx.Rollback() // read-only
	return s.getChunks(tx, req.Key, stream)
}

func (s *server) Set(stream pb.Mudah_SetServer) error {
	glog.Info("Got SET request")
	tx, err := s.db.Begin(true)
	if err != nil {
		return fmt.Errorf("failed to initiate tx: %v", err)
	}
	var ω int
	var β string
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to set value: %v", err)
		}
		β = chunk.Key
		err = tx.Bucket(keyBucketName).Put([]byte(fmt.Sprintf("%d~%s", ω, β)), chunk.Value)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to save chunk %d: %v", ω, err)
		}
		ω++
	}
	err = tx.Bucket(metaBucketName).Put([]byte(β), []byte(strconv.Itoa(ω)))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set metadata for %s: %v", β, err)
	}
	defer tx.Commit()
	glog.Infof("SET %s", β)
	return stream.SendAndClose(&pb.Key{Key: β})
}

func (s *server) List(req *pb.ListRequest, stream pb.Mudah_ListServer) error {
	glog.Infof("LIST %s", req.Prefix)
	tx, err := s.db.Begin(false)
	if err != nil {
		return fmt.Errorf("failed to initiate tx: %v", err)
	}
	defer tx.Rollback() // read-only

	c := tx.Bucket(metaBucketName).Cursor()
	var keys []string
	bytePrefix := []byte(req.Prefix)
	for ck, _ := c.Seek(bytePrefix); ck != nil && bytes.HasPrefix(ck, bytePrefix); ck, _ = c.Next() {
		keys = append(keys, string(ck))
	}
	for _, k := range keys {
		if err := s.getChunks(tx, k, stream); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	flag.Parse()
	if *dbPath == "" || *portNum <= 0 {
		flag.Usage()
		os.Exit(1)
	}
	s, err := newServer(*dbPath)
	if err != nil {
		glog.Fatalf("failed to create server: %v", err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *portNum))
	if err != nil {
		glog.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	pb.RegisterMudahServer(gs, s)
	gs.Serve(lis)
}
