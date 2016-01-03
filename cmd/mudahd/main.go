package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/grpclog/glogger"

	"golang.org/x/net/context"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/uluyol/mudahkv/pb"
)

var (
	keyBucketName = []byte("key")

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
	_, err = tx.CreateBucketIfNotExists(keyBucketName)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &server{db}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
	glog.Infof("GET %s", req.Key)
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate tx: %v", err)
	}
	defer tx.Rollback() // read-only
	value := tx.Bucket(keyBucketName).Get([]byte(req.Key))
	if value == nil {
		return nil, errors.New("no such file for key")
	}
	return &pb.GetReply{Key: req.Key, Value: string(value)}, nil
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetReply, error) {
	glog.Infof("SET %s", req.Key)
	tx, err := s.db.Begin(true)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate tx: %v", err)
	}
	err = tx.Bucket(keyBucketName).Put([]byte(req.Key), []byte(req.Value))
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to set value: %v", err)
	}
	defer tx.Commit()
	return &pb.SetReply{Key: req.Key, Value: req.Value}, nil
}

func (s *server) List(ctx context.Context, req *pb.ListRequest) (*pb.ListReply, error) {
	glog.Infof("LIST %s", req.Prefix)
	tx, err := s.db.Begin(false)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate tx: %v", err)
	}
	defer tx.Rollback() // read-only
	c := tx.Bucket(keyBucketName).Cursor()
	var values []*pb.GetReply
	bytePrefix := []byte(req.Prefix)
	for ck, cv := c.Seek(bytePrefix); ck != nil && bytes.HasPrefix(ck, bytePrefix); ck, cv = c.Next() {
		select {
		case <-ctx.Done():
			return nil, errors.New("deadline expired")
		default:
			values = append(values, &pb.GetReply{Key: string(ck), Value: string(cv)})
		}
	}
	return &pb.ListReply{Values: values}, nil
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
	pb.RegisterMudahKVServer(gs, s)
	gs.Serve(lis)
}
