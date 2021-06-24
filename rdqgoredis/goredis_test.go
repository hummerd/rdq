package rdqgoredis_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hummerd/rdq"
	"github.com/hummerd/rdq/rdqgoredis"
)

func Test_OneEventGet(t *testing.T) {
	raddr := os.Getenv("TEST_REDIS_HOST")
	if raddr == "" {
		t.Skip("no redis environment")
	}

	c := redis.NewClient(&redis.Options{
		Addr: raddr,
	})

	q := rdqgoredis.NewRDQ(&rdq.RDQOptions{
		Queue: "test-goredis-event-get",
	}, c)
	testData := "hi, it's me!"
	err := q.Add(context.Background(), time.Now().Add(time.Second*2), []byte(testData))
	if err != nil {
		t.Fatal(err)
	}

	_, d, err := q.Pop(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != testData {
		t.Fatal("unexpected response", string(d))
	}
}

func Test_WaitNoEvents(t *testing.T) {
	raddr := os.Getenv("TEST_REDIS_HOST")
	if raddr == "" {
		t.Skip("no redis environment")
	}

	c := redis.NewClient(&redis.Options{
		Addr: raddr,
	})

	q := rdqgoredis.NewRDQ(&rdq.RDQOptions{
		Queue: "test-goredis-no-event",
	}, c)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, _, err := q.Pop(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("unexpected error", err)
	}
}

func Test_CtxCancel(t *testing.T) {
	raddr := os.Getenv("TEST_REDIS_HOST")
	if raddr == "" {
		t.Skip("no redis environment")
	}

	c := redis.NewClient(&redis.Options{
		Addr: raddr,
	})

	q := rdqgoredis.NewRDQ(&rdq.RDQOptions{
		Queue: "test-goredis-ctx-cancel",
	}, c)

	testData := "hi, it's me!"
	err := q.Add(context.Background(), time.Now().Add(time.Second*10), []byte(testData))
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	s := time.Now()

	_, _, err = q.Pop(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("unexpected error", err)
	}

	if time.Since(s) > time.Second {
		t.Fatal("too long context cancelation", err)
	}
}
