// Topic: redis-basics | Concept: redis | Stage: 05-integrations
// Run:   go run ./05-integrations/redis/redis-basics
// Needs: docker compose -f 05-integrations/compose.yaml up -d --wait
//
// Explore:
//   - Redis is an in-memory key/value store; values can be strings, lists, hashes, sets
//   - Prefix keys by app and purpose ("go:basics:...") so apps sharing a Redis don't collide
//   - SET with an expiry (TTL) is how caches and sessions expire on their own
//   - redis.Nil is the "key does not exist" error: check it with errors.Is
//   - INCR is atomic: safe counters even with many clients
//   - Lists: queues and recent-items; Hashes: small objects (field -> value)
//   - Pipelines send many commands in one round trip
//   - Never use KEYS in production (it blocks Redis); use SCAN
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const hint = "Redis is not reachable. Start it with: docker compose -f 05-integrations/compose.yaml up -d --wait"
const prefix = "go:basics:"

// quietLogger hides go-redis's internal retry logs; the lesson prints its own hint.
type quietLogger struct{}

func (quietLogger) Printf(context.Context, string, ...any) {}

func redisURL() string {
	if u := os.Getenv("REDIS_URL"); u != "" {
		return u
	}
	return "redis://127.0.0.1:56379/0"
}

func must(err error, what string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", what, err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Topic: Redis basics with go-redis")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	redis.SetLogger(quietLogger{})
	opts, err := redis.ParseURL(redisURL())
	must(err, "parse REDIS_URL")
	opts.DialTimeout = 2 * time.Second
	rdb := redis.NewClient(opts)
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println(hint)
		os.Exit(1)
	}
	defer func() { // clean up every key this lesson created, using SCAN
		iter := rdb.Scan(context.Background(), 0, prefix+"*", 100).Iterator()
		n := 0
		for iter.Next(context.Background()) {
			rdb.Del(context.Background(), iter.Val())
			n++
		}
		fmt.Println("cleanup: deleted", n, "keys under", prefix)
	}()
	failed := false

	// Strings
	must(rdb.Set(ctx, prefix+"greeting", "hello", 0).Err(), "set")
	val, err := rdb.Get(ctx, prefix+"greeting").Result()
	fmt.Println("GET greeting:", val, err)

	_, err = rdb.Get(ctx, prefix+"missing").Result()
	fmt.Println("GET missing -> redis.Nil:", errors.Is(err, redis.Nil))

	// Expiry
	must(rdb.Set(ctx, prefix+"session", "user-42", 2*time.Second).Err(), "set ttl")
	ttl, _ := rdb.TTL(ctx, prefix+"session").Result()
	fmt.Println("session TTL:", ttl.Round(time.Second))

	// Atomic counter
	for range 5 {
		rdb.Incr(ctx, prefix+"page_views")
	}
	views, _ := rdb.IncrBy(ctx, prefix+"page_views", 10).Result()
	fmt.Println("page_views after 5 INCR + INCRBY 10:", views)
	if views != 15 {
		failed = true
	}

	// List as a queue: push on the right, pop from the left
	rdb.RPush(ctx, prefix+"jobs", "resize-image", "send-email", "build-report")
	jobs, _ := rdb.LRange(ctx, prefix+"jobs", 0, -1).Result()
	next, _ := rdb.LPop(ctx, prefix+"jobs").Result()
	left, _ := rdb.LLen(ctx, prefix+"jobs").Result()
	fmt.Println("jobs:", jobs, "| popped:", next, "| left:", left)

	// Hash as a small object
	rdb.HSet(ctx, prefix+"user:42", map[string]any{"name": "Asha", "city": "Pune", "plan": "pro"})
	user, _ := rdb.HGetAll(ctx, prefix+"user:42").Result()
	city, _ := rdb.HGet(ctx, prefix+"user:42", "city").Result()
	fmt.Println("user:42 hash:", user, "| city:", city)

	// Set: unique members
	rdb.SAdd(ctx, prefix+"tags", "go", "redis", "go", "docker")
	tags, _ := rdb.SCard(ctx, prefix+"tags").Result()
	fmt.Println("unique tags:", tags)
	if tags != 3 {
		failed = true
	}

	// Pipeline: 3 commands, 1 round trip
	pipe := rdb.Pipeline()
	a := pipe.Incr(ctx, prefix+"pipe")
	b := pipe.Incr(ctx, prefix+"pipe")
	c := pipe.Get(ctx, prefix+"greeting")
	_, err = pipe.Exec(ctx)
	must(err, "pipeline")
	fmt.Println("pipeline results:", a.Val(), b.Val(), c.Val())

	if failed {
		fmt.Println("MISMATCH")
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
