package rdqredigo

import (
	"context"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/hummerd/rdq"
)

func NewRDQ(opts *rdq.RDQOptions, redis redis.Conn) *rdq.RedisDelayedQueue {
	if redis == nil {
		panic("redis can not be nil")
	}

	opts.Redis = &Redis{redis}

	return rdq.NewRDQ(opts)
}

type Redis struct {
	client redis.Conn
}

func (r *Redis) BZPOPMIN(ctx context.Context, key string, timeout time.Duration) (float64, []byte, error) {
	res, err := r.client.Do("BZPOPMIN", key, timeout.Seconds())
	vals, err := redis.Values(res, err)
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return 0, nil, rdq.ErrNil
		}

		return 0, nil, err
	}

	var rkey string
	var data []byte
	var score float64
	_, err = redis.Scan(vals, &rkey, &data, &score)
	if err != nil {
		return 0, nil, err
	}

	return score, data, nil
}

func (r *Redis) ZADD(ctx context.Context, key string, score float64, memeber []byte) error {
	_, err := r.client.Do("ZADD", key, score, memeber)
	return err
}
