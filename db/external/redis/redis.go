package redis

import (
	"context"
	"strconv"

	"github.com/snowmerak/twisted-lyfes/db/external"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	rdb *redis.Client
	ctx context.Context
}

func Connect(url string, port int, password string) external.KVCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url + ":" + strconv.Itoa(port),
		Password: password,
		DB:       0,
	})

	return &RedisDB{rdb: rdb, ctx: context.Background()}
}

func (r *RedisDB) Set(id, url string) error {
	if err := r.rdb.SetNX(r.ctx, id, url, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) Get(id string) (string, error) {
	url, err := r.rdb.Get(r.ctx, id).Result()
	if err != nil {
		return "", err
	}
	return url, nil
}
