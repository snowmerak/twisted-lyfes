package badger

import (
	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"time"
)

type DB struct {
	*badger.DB
}

func (d *DB) Put(subject string, key string, data protoreflect.ProtoMessage) error {
	value, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	if err := d.DB.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(subject+key), value); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (d *DB) Get(subject string, key string, value protoreflect.ProtoMessage) (protoreflect.ProtoMessage, error) {
	var data []byte = nil
	if err := d.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(subject + key))
		if err != nil {
			return err
		}
		if err := item.Value(func(val []byte) error {
			data = val
			return nil
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(data, value); err != nil {
		return nil, err
	}
	return value, nil
}

func (d *DB) PutWithExpire(subject string, ttl int64, key string, data protoreflect.ProtoMessage) error {
	value, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	if err := d.DB.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(subject+key), value).WithTTL(time.Duration(ttl))
		if err := txn.SetEntry(entry); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (d *DB) Delete(subject string, key string) error {
	if err := d.DB.Update(func(txn *badger.Txn) error {
		if err := txn.Delete([]byte(subject + key)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (d *DB) Exists(subject string, key string) (bool, error) {
	var exists bool = false
	if err := d.DB.View(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte(subject + key)); err != nil {
			if err == badger.ErrKeyNotFound {
				exists = false
				return nil
			}
			return err
		}
		exists = true
		return nil
	}); err != nil {
		return false, err
	}
	return exists, nil
}

func (d *DB) Close() error {
	return d.DB.Close()
}
