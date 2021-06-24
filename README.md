# Redis Delayed Queue (RDQ)
RDQ is the queue of delayed events persisted in Redis written in GO.

## Examples

### Redigo example
```go
import "github.com/gomodule/redigo/redis"

func AddAndGet() {
    rc, _ := redis.Dial("tcp", redisAddress)

    q := rdqredigo.NewRDQ(&rdq.RDQOptions{
        Queue: "test-queue",
    }, rc)

    testData := "hi, it's me!"
    _ = q.AddAfter(context.Background(), time.Second*10, []byte(testData))

    // Pop will return after 10 seconds
    eventTime, data, err = q.Pop(context.Background())
}
```

### GoRedis example
```go
import "github.com/go-redis/redis/v8"

func AddAndGet() {
    c := redis.NewClient(&redis.Options{
        Addr: raddr,
    })

    q := rdqgoredis.NewRDQ(&rdq.RDQOptions{
        Queue: "test-queue",
    }, c)

    testData := "hi, it's me!"
    _ = q.AddAfter(context.Background(), time.Second*10, []byte(testData))

    // Pop will return after 10 seconds
    eventTime, data, err = q.Pop(context.Background())
}
```
