package rdqredigo

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/hummerd/rdq"
)

func NewRDQPool(opts *rdq.RDQOptions, redisPool *redis.Pool) *rdq.RedisDelayedQueue {
	if redisPool == nil {
		panic("redisPool can not be nil")
	}

	opts.Redis = &RedisPool{redisPool}

	return rdq.NewRDQ(opts)
}

type RedisPool struct {
	redisPool *redis.Pool
}

func (rp *RedisPool) BZPOPMIN(ctx context.Context, key string, timeout time.Duration) (float64, []byte, error) {
	c := rp.redisPool.Get()
	defer c.Close()

	r := Redis{c}
	return r.BZPOPMIN(ctx, key, timeout)
}

func (rp *RedisPool) ZADD(ctx context.Context, key string, score float64, memeber []byte) error {
	c := rp.redisPool.Get()
	defer c.Close()

	r := Redis{c}
	return r.ZADD(ctx, key, score, memeber)
}
