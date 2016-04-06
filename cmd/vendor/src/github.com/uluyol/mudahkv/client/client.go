package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/uluyol/mudahkv/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	_MB       = 1 << 20
	chunkSize = 8 * _MB
)

type Client struct {
	cc *grpc.ClientConn
	mc pb.MudahClient
}

func Dial(addr string) (*Client, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{cc, pb.NewMudahClient(cc)}, nil
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	r, err := c.GetStream(ctx, key)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func (c *Client) GetStream(ctx context.Context, key string) (io.Reader, error) {
	stream, err := c.mc.Get(ctx, &pb.Key{Key: key})
	if err != nil {
		return nil, err
	}
	iter := newReaderIter(stream)
	if !iter.Next() {
		return nil, fmt.Errorf("unable to find key: %v", iter.Err())
	}
	return iter.Value(), nil
}

func (c *Client) Set(ctx context.Context, key string, value []byte) error {
	if len(value) > chunkSize {
		return c.SetStream(ctx, key, bytes.NewReader(value))
	}
	stream, err := c.mc.Set(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&pb.KVChunk{Key: key, Value: value}); err != nil {
		return err
	}
	_, err = stream.CloseAndRecv()
	return err
}

func (c *Client) SetStream(ctx context.Context, key string, r io.Reader) error {
	stream, err := c.mc.Set(ctx)
	if err != nil {
		return err
	}
	var chunk [chunkSize]byte
	for {
		n, rerr := r.Read(chunk[:])
		if rerr != nil && rerr != io.EOF {
			return rerr
		}
		if err := stream.Send(&pb.KVChunk{Key: key, Value: chunk[:n]}); err != nil {
			return err
		}
		if rerr == io.EOF {
			break
		}
	}
	_, err = stream.CloseAndRecv()
	return err
}

func (c *Client) List(ctx context.Context, prefix string) ([]KV, error) {
	iter, err := c.ListStream(ctx, prefix)
	if err != nil {
		return nil, err
	}
	var kvs []KV
	for iter.Next() {
		v, err := ioutil.ReadAll(iter.Value())
		if err != nil {
			return kvs, fmt.Errorf("error while reading data: %v", err)
		}
		kvs = append(kvs, KV{iter.Key(), v})
	}
	return kvs, iter.Err()
}

func (c *Client) ListStream(ctx context.Context, prefix string) (ReaderIter, error) {
	stream, err := c.mc.List(ctx, &pb.ListRequest{Prefix: prefix})
	if err != nil {
		return nil, err
	}
	return newReaderIter(stream), nil
}

type ReaderIter interface {
	Next() bool
	Key() string
	Value() io.Reader
	Err() error
}

type kvChunkReceiver interface {
	Recv() (*pb.KVChunk, error)
}

func newReaderIter(stream kvChunkReceiver) ReaderIter {
	return &readerIter{first: true, stream: stream}
}

type readerIter struct {
	stream kvChunkReceiver
	first  bool
	k      string
	last   *pb.KVChunk
	err    error
}

func (iter *readerIter) Next() bool {
	if iter.err != nil {
		return false
	}
	if iter.last != nil {
		iter.k = iter.last.Key
		return true
	}
	for {
		chunk, err := iter.stream.Recv()
		if err == io.EOF {
			return false
		} else if err != nil {
			iter.err = err
			return false
		}
		if !iter.first && iter.last.Key == iter.k {
			continue
		}
		iter.last = chunk
		iter.k = chunk.Key
		iter.first = false
		break
	}
	return true
}

func (iter *readerIter) Key() string {
	return iter.k
}

func (iter *readerIter) Value() io.Reader {
	if iter.err != nil {
		return nil
	}
	return reader{iter}
}

func (iter *readerIter) Err() error {
	return iter.err
}

type reader struct {
	p *readerIter
}

func (r reader) Read(b []byte) (int, error) {
	if r.p.last == nil || r.p.last.Key != r.p.k {
		return 0, io.EOF
	}
	var n int
	for n < len(b) {
		m := copy(b, r.p.last.Value)
		n += m
		b = b[m:]
		r.p.last.Value = r.p.last.Value[m:]
		if n < len(b) {
			chunk, err := r.p.stream.Recv()
			if err == io.EOF {
				r.p.last = nil
				return n, io.EOF
			} else if err != nil {
				r.p.err = err
				return n, err
			}
			r.p.last = chunk
			if chunk.Key != r.p.k {
				return n, io.EOF
			}
		}
	}
	return n, nil
}

type KV struct {
	Key   string
	Value []byte
}

func (c *Client) Close() error {
	return c.cc.Close()
}
