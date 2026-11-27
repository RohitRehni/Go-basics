// Topic: http-client | Concept: http | Stage: 05-integrations
// Run:   go run ./05-integrations/http/http-client
//
// Explore:
//   - Never use http.Get / http.DefaultClient in real code: they have NO timeout
//   - Build one http.Client with Timeout and reuse it (it pools connections)
//   - Per-request deadlines: http.NewRequestWithContext + context.WithTimeout
//   - Always `defer resp.Body.Close()`, and check resp.StatusCode: a 404 or 500
//     is NOT a Go error
//   - Decode JSON straight from the body with json.NewDecoder
//   - Retry only idempotent requests, only on 5xx / network errors, with backoff
//   - httptest.NewServer starts a real server on 127.0.0.1 for tests and demos:
//     this lesson never talks to the internet
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// newTestServer plays the role of a remote API.
func newTestServer() *httptest.Server {
	var flakyCalls atomic.Int32
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("id") != "1" {
			http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{1, "Asha", "asha@example.com"})
	})
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var u User
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil || u.Name == "" {
			http.Error(w, `{"error":"name is required"}`, http.StatusBadRequest)
			return
		}
		u.ID = 2
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	})
	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
			fmt.Fprintln(w, "finally")
		case <-r.Context().Done(): // client gave up
		}
	})
	mux.HandleFunc("GET /flaky", func(w http.ResponseWriter, r *http.Request) {
		if flakyCalls.Add(1) <= 2 {
			http.Error(w, "try again", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "ok on attempt", flakyCalls.Load())
	})
	return httptest.NewServer(mux)
}

// getJSON GETs url and decodes a JSON body into out, treating non-2xx as an error.
func getJSON(ctx context.Context, client *http.Client, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("GET %s: status %d: %s", url, resp.StatusCode, bytes.TrimSpace(body))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

// getWithRetry retries 5xx responses with exponential backoff.
func getWithRetry(ctx context.Context, client *http.Client, url string, attempts int) (string, error) {
	backoff := 50 * time.Millisecond
	for i := 1; i <= attempts; i++ {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		resp, err := client.Do(req)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode < 500 {
				return string(bytes.TrimSpace(body)), nil
			}
			err = fmt.Errorf("status %d", resp.StatusCode)
		}
		fmt.Printf("    attempt %d failed (%v), waiting %v\n", i, err, backoff)
		time.Sleep(backoff)
		backoff *= 2
	}
	return "", fmt.Errorf("gave up after %d attempts", attempts)
}

func main() {
	fmt.Println("Topic: HTTP client")
	srv := newTestServer()
	defer srv.Close()
	fmt.Println("test API running at", srv.URL)

	client := &http.Client{Timeout: 5 * time.Second}
	ctx := context.Background()
	failed := false

	// GET + JSON decode
	var u User
	err := getJSON(ctx, client, srv.URL+"/users/1", &u)
	fmt.Printf("GET /users/1  -> %+v err=%v\n", u, err)
	if err != nil || u.Name != "Asha" {
		failed = true
	}

	// 404 is a status code, not a transport error
	err = getJSON(ctx, client, srv.URL+"/users/99", &u)
	fmt.Println("GET /users/99 ->", err)
	if err == nil {
		failed = true
	}

	// POST JSON
	payload, _ := json.Marshal(User{Name: "Ravi", Email: "ravi@example.com"})
	resp, err := client.Post(srv.URL+"/users", "application/json", bytes.NewReader(payload))
	if err != nil {
		fmt.Fprintln(os.Stderr, "POST failed:", err)
		os.Exit(1)
	}
	var created User
	json.NewDecoder(resp.Body).Decode(&created)
	resp.Body.Close()
	fmt.Printf("POST /users   -> %d %+v\n", resp.StatusCode, created)
	if resp.StatusCode != http.StatusCreated || created.ID != 2 {
		failed = true
	}

	resp, _ = client.Post(srv.URL+"/users", "application/json", bytes.NewReader([]byte(`{}`)))
	resp.Body.Close()
	fmt.Println("POST /users with no name ->", resp.StatusCode)

	// Per-request timeout shorter than the client's
	reqCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	start := time.Now()
	req, _ := http.NewRequestWithContext(reqCtx, http.MethodGet, srv.URL+"/slow", nil)
	_, err = client.Do(req)
	cancel()
	fmt.Printf("GET /slow with 200ms deadline -> gave up after %v, deadline exceeded: %v\n",
		time.Since(start).Round(10*time.Millisecond), errors.Is(err, context.DeadlineExceeded))
	if !errors.Is(err, context.DeadlineExceeded) {
		failed = true
	}

	// Retry with backoff
	fmt.Println("GET /flaky with retries:")
	body, err := getWithRetry(ctx, client, srv.URL+"/flaky", 4)
	fmt.Println("    result:", body, err)
	if err != nil {
		failed = true
	}

	if failed {
		fmt.Println("MISMATCH")
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
