package redis

import (
	"context"
	"github.com/rueian/rueidis"
)

type Client struct {
	client rueidis.Client

	ctx    context.Context
	cancel context.CancelFunc
}

func NewClient(ctx context.Context, addr ...string) (*Client, error) {
	c, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: addr,
	})
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	return &Client{
		client: c,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (c *Client) Close() {
	c.cancel()
	c.client.Close()
}
