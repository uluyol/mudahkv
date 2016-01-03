package client

import (
	"github.com/uluyol/mudahkv/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type KV struct {
	Key   string
	Value string
}

type Client struct {
	cc *grpc.ClientConn
	mc pb.MudahKVClient
}

func Dial(addr string) (*Client, error) {
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{cc, pb.NewMudahKVClient(cc)}, nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	resp, err := c.mc.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		return "", err
	}
	return resp.Value, err
}

func (c *Client) Set(ctx context.Context, key, value string) error {
	_, err := c.mc.Set(ctx, &pb.SetRequest{Key: key, Value: value})
	return err
}

func (c *Client) List(ctx context.Context, prefix string) ([]KV, error) {
	resp, err := c.mc.List(ctx, &pb.ListRequest{Prefix: prefix})
	if err != nil {
		return nil, err
	}
	kv := make([]KV, len(resp.Values))
	for i, v := range resp.Values {
		kv[i] = KV{v.Key, v.Value}
	}
	return kv, nil
}

func (c *Client) Close() error {
	return c.cc.Close()
}
