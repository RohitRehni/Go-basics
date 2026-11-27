package main

// cache.go: the caching layer. Only this file knows Redis.

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const cachePrefix = "go:notes:"

type NoteCache struct {
	rdb *redis.Client
	ttl time.Duration
}

// quietLogger hides go-redis's internal retry logs; the lesson prints its own hint.
type quietLogger struct{}

func (quietLogger) Printf(context.Context, string, ...any) {}

func NewNoteCache(ctx context.Context, url string, ttl time.Duration) (*NoteCache, error) {
	redis.SetLogger(quietLogger{})
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	opts.DialTimeout = 2 * time.Second
	rdb := redis.NewClient(opts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, err
	}
	return &NoteCache{rdb: rdb, ttl: ttl}, nil
}

func (c *NoteCache) Close() error { return c.rdb.Close() }

func (c *NoteCache) key(id int) string { return fmt.Sprintf("%s%d", cachePrefix, id) }

// Get returns the cached note and true on a hit. Any Redis problem counts as a miss.
func (c *NoteCache) Get(ctx context.Context, id int) (Note, bool) {
	raw, err := c.rdb.Get(ctx, c.key(id)).Bytes()
	if err != nil {
		return Note{}, false
	}
	var n Note
	if json.Unmarshal(raw, &n) != nil {
		return Note{}, false
	}
	return n, true
}

func (c *NoteCache) Set(ctx context.Context, n Note) {
	if data, err := json.Marshal(n); err == nil {
		c.rdb.Set(ctx, c.key(n.ID), data, c.ttl) // errors ignored: the cache is optional
	}
}

func (c *NoteCache) Invalidate(ctx context.Context, id int) {
	c.rdb.Del(ctx, c.key(id))
}

// Clear deletes every key this app owns (SCAN, never KEYS) and returns how many.
func (c *NoteCache) Clear(ctx context.Context) int {
	n := 0
	iter := c.rdb.Scan(ctx, 0, cachePrefix+"*", 100).Iterator()
	for iter.Next(ctx) {
		c.rdb.Del(ctx, iter.Val())
		n++
	}
	return n
}
