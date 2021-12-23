package mq

type MQ interface {
	Publish(subject string, data []byte) error
	Subscribe(subject string, channel chan<- []byte) error
	SetReply(subject string, handler func([]byte) []byte) error
	Unsubscribe(subject string) error
	Disconnect() error
}
