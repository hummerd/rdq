package rdq

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// RDQOptions is a settings for RDQ.
type RDQOptions struct {
	// Queue is the name of the ZSet in redis
	Queue string
	// Redis is redis client
	Redis Redis
	// Now is function returning current time (usefull for tests). By default time.Now.
	Now func() time.Time
	// PollInterval to check if there is appropriate events. By default 1 second.
	PollInterval time.Duration
	// ReturnTimeout is timeout for returning event that is not came yet. By default 100 msec.
	ReturnTimeout time.Duration
}

// NewRDQ creates new RedisDelayedQueue
func NewRDQ(opts *RDQOptions) *RedisDelayedQueue {
	if opts.Queue == "" {
		panic("queue can not be empty")
	}

	if opts.Redis == nil {
		panic("redis can not be nil")
	}

	if opts.Now == nil {
		opts.Now = time.Now
	}

	if opts.PollInterval <= 0 {
		opts.PollInterval = time.Second
	}

	if opts.ReturnTimeout <= 0 {
		opts.ReturnTimeout = time.Millisecond * 100
	}

	return &RedisDelayedQueue{
		opts: *opts,
	}
}

var ErrNil = errors.New("redis has no data")

type Redis interface {
	BZPOPMIN(ctx context.Context, key string, timeout time.Duration) (float64, []byte, error)
	ZADD(ctx context.Context, key string, score float64, memeber []byte) error
}

// RedisDelayedQueue is delayed queue persisted in redis.
// You can call Add from multiple goroutines, but you should call to Pop in single goroutine.
type RedisDelayedQueue struct {
	opts      RDQOptions
	pollTimer *time.Timer
}

// AddAfter adds new delayed event to queue.
func (rdq *RedisDelayedQueue) AddAfter(ctx context.Context, delay time.Duration, item []byte) error {
	return rdq.Add(ctx, rdq.opts.Now().Add(delay), item)
}

// Add adds new delayed event to queue.
func (rdq *RedisDelayedQueue) Add(ctx context.Context, at time.Time, item []byte) error {
	return rdq.opts.Redis.ZADD(ctx, rdq.opts.Queue, float64(at.UnixNano()), item)
}

// Pop blocks until get event with appropriate time. Pop is not safe for concurrent use.
func (rdq *RedisDelayedQueue) Pop(ctx context.Context) (time.Time, []byte, error) {
	for {
		s, d, err := rdq.opts.Redis.BZPOPMIN(ctx, rdq.opts.Queue, rdq.opts.PollInterval)
		if err != nil {
			if errors.Is(err, ErrNil) {
				if ctx.Err() == nil {
					continue
				}

				return time.Time{}, nil, fmt.Errorf("error on waiting event: %w", ctx.Err())
			}
			return time.Time{}, nil, fmt.Errorf("error on waiting event: %w", err)
		}

		tt := time.Unix(0, int64(s))
		if tt.After(rdq.opts.Now()) {
			rctx, rcancel := context.WithTimeout(context.Background(), rdq.opts.ReturnTimeout)
			err = rdq.Add(rctx, tt, d)
			rcancel()

			if err != nil {
				return time.Time{}, nil, fmt.Errorf("error on returning event: %w", err)
			}

			if rdq.pollTimer == nil {
				rdq.pollTimer = time.NewTimer(rdq.opts.PollInterval)
			} else {
				rdq.pollTimer.Reset(rdq.opts.PollInterval)
			}

			select {
			case <-ctx.Done():
				if !rdq.pollTimer.Stop() {
					<-rdq.pollTimer.C
				}
				return time.Time{}, nil, ctx.Err()
			case <-rdq.pollTimer.C:
				rdq.pollTimer.Stop()
				continue
			}
		}

		return tt, d, nil
	}
}
