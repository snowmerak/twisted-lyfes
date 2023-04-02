package nats

import (
	"encoding/hex"
	nats "github.com/nats-io/nats.go"
	"golang.org/x/exp/slog"
)

type Connection struct {
	nc     *nats.Conn
	subs   map[string]*nats.Subscription
	logger *slog.Logger
}

func Connect(url string, logger *slog.Logger) (*Connection, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Connection{nc, make(map[string]*nats.Subscription), logger}, nil
}

func (c *Connection) Publish(subject string, data []byte) error {
	return c.nc.Publish(subject, data)
}

func (c *Connection) Subscribe(subject string, handler func([]byte) []byte) error {
	sub, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
		resp := handler(msg.Data)
		if msg.Reply == "" {
			return
		}
		if err := c.nc.Publish(msg.Reply, resp); err != nil {
			c.logger.Error("failed to reply",
				slog.String("subject", subject),
				slog.String("reply_to", msg.Reply),
				slog.String("contents", hex.EncodeToString(resp)))
		}
	})
	if err != nil {
		return err
	}
	c.subs[subject] = sub
	return nil
}

func (c *Connection) SubscribeQueue(queue string, subject string, handler func([]byte) []byte) error {
	sub, err := c.nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		resp := handler(msg.Data)
		if msg.Reply == "" {
			return
		}
		if err := c.nc.Publish(msg.Reply, resp); err != nil {
			c.logger.Error("failed to reply",
				slog.String("subject", subject),
				slog.String("queue", queue),
				slog.String("reply_to", msg.Reply),
				slog.String("contents", hex.EncodeToString(resp)))
		}
	})
	if err != nil {
		return err
	}
	c.subs[subject] = sub
	return nil
}

func (c *Connection) Unsubscribe(subject string) error {
	sub, ok := c.subs[subject]
	if !ok {
		return nil
	}
	if err := sub.Unsubscribe(); err != nil {
		return err
	}
	delete(c.subs, subject)
	return nil
}

func (c *Connection) Disconnect() error {
	for _, sub := range c.subs {
		if err := sub.Unsubscribe(); err != nil {
			return err
		}
	}
	if err := c.nc.Drain(); err != nil {
		return err
	}
	c.nc.Close()
	return nil
}
