package db

import (
	"github.com/twisted-lyfes/certificate-log/model"
	"github.com/xujiajun/nutsdb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type DB struct {
	db *nutsdb.DB
}

func New(path string) (*DB, error) {
	opt := nutsdb.DefaultOptions
	opt.Dir = path
	db, err := nutsdb.Open(opt)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Put(subject string, key string, data protoreflect.ProtoMessage) error {
	bs, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	return db.db.Update(func(tx *nutsdb.Tx) error {
		err := tx.Put(subject, []byte(key), bs, 0)
		return err
	})
}

func (db *DB) PutWithExpire(subject string, ttl int64, timestamp int64, key string, data protoreflect.ProtoMessage) error {
	bs, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	return db.db.Update(func(tx *nutsdb.Tx) error {
		err := tx.PutWithTimestamp(subject, []byte(key), bs, uint32(ttl), uint64(timestamp))
		return err
	})
}

func (db *DB) Get(subject string, key string) (protoreflect.ProtoMessage, error) {
	var value []byte = nil
	err := db.db.View(func(tx *nutsdb.Tx) error {
		entry, err := tx.Get(subject, []byte(key))
		if err != nil {
			return err
		}
		value = entry.Value
		return nil
	})
	if err != nil {
		return nil, err
	}
	data := new(model.CertificateLog)
	proto.Unmarshal(value, data)
	return data, nil
}
