package rdqgoredis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hummerd/rdq"
)

func NewRDQ(opts *rdq.RDQOptions, redis redis.Cmdable) *rdq.RedisDelayedQueue {
	if redis == nil {
		panic("redis can not be nil")
	}

	opts.Redis = &Redis{redis}

	return rdq.NewRDQ(opts)
}

type Redis struct {
	client redis.Cmdable
}

func (r *Redis) BZPOPMIN(ctx context.Context, key string, timeout time.Duration) (float64, []byte, error) {
	res, err := r.client.BZPopMin(ctx, timeout, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, nil, rdq.ErrNil
		}

		return 0, nil, err
	}

	var d []byte
	switch v := res.Member.(type) {
	case string:
		d = []byte(v)
	case []byte:
		d = v
	}

	return res.Score, d, nil
}

func (r *Redis) ZADD(ctx context.Context, key string, score float64, memeber []byte) error {
	_, err := r.client.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: memeber,
	}).Result()
	return err
}
