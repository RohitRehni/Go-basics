// Topic: redis-cache | Concept: redis | Stage: 05-integrations
// Run:   go run ./05-integrations/redis/redis-cache
// Needs: docker compose -f 05-integrations/compose.yaml up -d --wait
//
// Explore:
//   - Cache-aside: read the cache; on a miss, load from the slow source and store it with a TTL
//   - The TTL bounds how stale data can get, even if you forget to invalidate
//   - On writes, delete (invalidate) the cached key so the next read reloads it
//   - Store structured values as JSON
//   - If Redis is down, a cache should degrade to the slow path, not fail the request
//   - Watch out for a "thundering herd": many misses at once all hit the slow source
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const hint = "Redis is not reachable. Start it with: docker compose -f 05-integrations/compose.yaml up -d --wait"
const prefix = "go:cache:"

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// The "database": slow on purpose.
var products = map[int]Product{
	1: {1, "Keyboard", 49.99},
	2: {2, "Mouse", 19.50},
}
var slowLoads int

func loadProduct(id int) (Product, error) {
	slowLoads++
	time.Sleep(200 * time.Millisecond)
	p, ok := products[id]
	if !ok {
		return Product{}, fmt.Errorf("product %d not found", id)
	}
	return p, nil
}

// quietLogger hides go-redis's internal retry logs; the lesson prints its own hint.
type quietLogger struct{}

func (quietLogger) Printf(context.Context, string, ...any) {}

type Cache struct {
	rdb *redis.Client
	ttl time.Duration
}

func (c *Cache) key(id int) string { return fmt.Sprintf("%sproduct:%d", prefix, id) }

// GetProduct implements cache-aside and reports whether it was a hit.
func (c *Cache) GetProduct(ctx context.Context, id int) (Product, bool, error) {
	raw, err := c.rdb.Get(ctx, c.key(id)).Bytes()
	if err == nil {
		var p Product
		if json.Unmarshal(raw, &p) == nil {
			return p, true, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		fmt.Println("  (cache error, falling back to the database:", err, ")")
	}

	p, err := loadProduct(id)
	if err != nil {
		return Product{}, false, err
	}
	if data, err := json.Marshal(p); err == nil {
		c.rdb.Set(ctx, c.key(id), data, c.ttl)
	}
	return p, false, nil
}

// UpdatePrice writes to the source of truth, then invalidates the cache.
func (c *Cache) UpdatePrice(ctx context.Context, id int, price float64) {
	p := products[id]
	p.Price = price
	products[id] = p
	c.rdb.Del(ctx, c.key(id))
}

func main() {
	fmt.Println("Topic: cache-aside with Redis")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	redis.SetLogger(quietLogger{})
	opts, err := redis.ParseURL(func() string {
		if u := os.Getenv("REDIS_URL"); u != "" {
			return u
		}
		return "redis://127.0.0.1:56379/0"
	}())
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad REDIS_URL:", err)
		os.Exit(1)
	}
	opts.DialTimeout = 2 * time.Second
	rdb := redis.NewClient(opts)
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Println(hint)
		os.Exit(1)
	}
	cache := &Cache{rdb: rdb, ttl: 500 * time.Millisecond}
	rdb.Del(ctx, cache.key(1), cache.key(2))
	defer func() {
		rdb.Del(context.Background(), cache.key(1), cache.key(2))
		fmt.Println("cleanup: deleted cache keys")
	}()

	get := func(label string, id int) bool {
		start := time.Now()
		p, hit, err := cache.GetProduct(ctx, id)
		status := "MISS"
		if hit {
			status = "HIT "
		}
		if err != nil {
			fmt.Printf("  %-26s error: %v\n", label, err)
			return false
		}
		fmt.Printf("  %-26s %s %-8s %6.2f  in %v\n", label, status, p.Name, p.Price, time.Since(start).Round(time.Millisecond))
		return hit
	}

	failed := false
	if get("first read of product 1", 1) {
		failed = true // must be a miss
	}
	if !get("second read of product 1", 1) {
		failed = true // must be a hit
	}
	get("first read of product 2", 2)

	cache.UpdatePrice(ctx, 1, 44.99)
	fmt.Println("  updated product 1 price -> cache key deleted")
	if get("read after update", 1) {
		failed = true // must reload
	}

	time.Sleep(600 * time.Millisecond) // longer than the 500ms TTL
	if get("read after TTL expired", 2) {
		failed = true // must reload
	}
	get("unknown product 9", 9)

	fmt.Println("slow loads from the 'database':", slowLoads)
	if slowLoads != 5 {
		failed = true
	}

	if failed {
		fmt.Println("MISMATCH in hit/miss pattern")
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
