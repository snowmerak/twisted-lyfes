package nats

import (
	nats "github.com/nats-io/nats.go"
	"github.com/snowmerak/twisted-lyfes/src/mq"
)

type Connection struct {
	nc   *nats.Conn
	subs map[string]*nats.Subscription
}

func Connect(url string) (mq.MQ, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Connection{nc, make(map[string]*nats.Subscription)}, nil
}

func (c *Connection) Publish(subject string, data []byte) error {
	return c.nc.Publish(subject, data)
}

func (c *Connection) Subscribe(subject string, channel chan<- []byte) error {
	sub, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
		channel <- msg.Data
	})
	if err != nil {
		return err
	}
	c.subs[subject] = sub
	return nil
}

func (c *Connection) SetReply(subject string, handler func([]byte) []byte) error {
	sub, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
		c.nc.Publish(msg.Reply, handler(msg.Data))
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
	sub.Unsubscribe()
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
