package embedded

import "google.golang.org/protobuf/reflect/protoreflect"

type EmbeddedDB interface {
	Put(subject string, key string, data protoreflect.ProtoMessage) error
	Get(subject string, key string, value protoreflect.ProtoMessage) (protoreflect.ProtoMessage, error)
	PutWithExpire(subject string, ttl int64, timestamp int64, key string, data protoreflect.ProtoMessage) error
}
