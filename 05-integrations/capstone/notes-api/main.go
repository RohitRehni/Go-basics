// Topic: notes-api | Concept: capstone | Stage: 05-integrations
// Run:   go run ./05-integrations/capstone/notes-api
// Needs: docker compose -f 05-integrations/compose.yaml up -d --wait
//
// Multi-file lesson (all files are package main):
//   - main.go      wires everything together, starts the server, plays the client
//   - store.go     NoteStore: PostgreSQL via a pgxpool connection pool
//   - cache.go     NoteCache: read-through cache in Redis with a TTL
//   - handlers.go  HTTP routes: JSON in, JSON out, X-Cache: HIT/MISS header
//
// Explore:
//   - Layers: handlers know HTTP, the store knows SQL, the cache knows Redis;
//     each can be tested or replaced on its own
//   - Read-through cache: GET checks Redis first, falls back to Postgres, then fills Redis
//   - Writes go to Postgres first, then invalidate the cached copy
//   - The cache is optional: if Redis fails, requests still work (just slower)
//   - Graceful shutdown: stop accepting requests, finish in-flight ones, close pools
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

const hint = "Postgres/Redis not reachable. Start them with: docker compose -f 05-integrations/compose.yaml up -d --wait"

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	fmt.Println("Topic: notes API (net/http + Postgres + Redis)")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	store, err := NewNoteStore(ctx, env("DATABASE_URL", "postgres://learning:learning@127.0.0.1:55432/learning"))
	if err != nil {
		fmt.Println(hint)
		os.Exit(1)
	}
	defer store.Close()

	cache, err := NewNoteCache(ctx, env("REDIS_URL", "redis://127.0.0.1:56379/0"), 30*time.Second)
	if err != nil {
		fmt.Println(hint)
		os.Exit(1)
	}
	defer cache.Close()

	if err := store.CreateSchema(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "schema:", err)
		os.Exit(1)
	}
	defer func() { // leave nothing behind, so the lesson can run again
		store.DropSchema(context.Background())
		n := cache.Clear(context.Background())
		fmt.Printf("cleanup: dropped go_notes, deleted %d cache keys\n", n)
	}()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	srv := &http.Server{Handler: NewAPI(store, cache).Routes(), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintln(os.Stderr, "serve:", err)
		}
	}()
	base := "http://" + ln.Addr().String()
	fmt.Println("notes API listening on", base)

	failed := !runClient(base)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintln(os.Stderr, "shutdown:", err)
		failed = true
	}
	fmt.Println("server shut down gracefully")

	if failed {
		fmt.Println("MISMATCH: an API call did not behave as expected")
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}

// runClient exercises every endpoint and reports whether all checks passed.
func runClient(base string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	ok := true

	do := func(method, path string, body any, wantStatus int, wantCache string) map[string]any {
		var buf bytes.Buffer
		if body != nil {
			json.NewEncoder(&buf).Encode(body)
		}
		req, _ := http.NewRequest(method, base+path, &buf)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("  request failed:", err)
			ok = false
			return nil
		}
		defer resp.Body.Close()
		var out map[string]any
		var raw bytes.Buffer
		raw.ReadFrom(resp.Body)
		json.Unmarshal(raw.Bytes(), &out)

		cacheHdr := resp.Header.Get("X-Cache")
		mark := "ok"
		if resp.StatusCode != wantStatus || (wantCache != "" && cacheHdr != wantCache) {
			mark = fmt.Sprintf("MISMATCH (want %d %s)", wantStatus, wantCache)
			ok = false
		}
		fmt.Printf("  %-6s %-9s %d %-4s %s  %s\n", method, path, resp.StatusCode, cacheHdr, mark, bytes.TrimSpace(raw.Bytes()))
		return out
	}

	created := do("POST", "/notes", map[string]string{"title": "Groceries", "body": "milk, eggs"}, 201, "")
	do("POST", "/notes", map[string]string{"title": "Ideas", "body": "web scraper for meter data"}, 201, "")
	do("POST", "/notes", map[string]string{"body": "no title"}, 400, "")
	id := fmt.Sprint(created["id"])
	do("GET", "/notes", nil, 200, "")
	do("GET", "/notes/"+id, nil, 200, "MISS") // first read: from Postgres, then cached
	do("GET", "/notes/"+id, nil, 200, "HIT")  // second read: straight from Redis
	do("PUT", "/notes/"+id, map[string]string{"title": "Groceries", "body": "milk, eggs, bread"}, 200, "")
	do("GET", "/notes/"+id, nil, 200, "MISS") // the update invalidated the cache
	do("GET", "/notes/"+id, nil, 200, "HIT")
	do("DELETE", "/notes/"+id, nil, 204, "")
	do("GET", "/notes/"+id, nil, 404, "")
	do("GET", "/notes/999", nil, 404, "")
	return ok
}
