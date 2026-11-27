// Topic: rest-api | Concept: http | Stage: 05-integrations
// Run:   go run ./05-integrations/http/rest-api
//
// Explore:
//   - Go 1.22+ ServeMux patterns: "GET /todos/{id}", read the value with r.PathValue("id")
//   - REST: GET list / GET one / POST create (201) / PUT update / DELETE (204)
//   - Status codes matter: 400 bad input, 404 missing, 405 wrong method (automatic)
//   - Handlers run concurrently: protect shared state with a mutex
//   - Middleware = a function that wraps a handler (logging, auth, recovery)
//   - Listen on "127.0.0.1:0" to get a free port; http.Server.Shutdown drains
//     in-flight requests before stopping
//   - In a real service, swap the in-memory map for a database (see capstone/notes-api)
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
	"slices"
	"strconv"
	"sync"
	"time"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type store struct {
	mu     sync.Mutex
	todos  map[int]Todo
	nextID int
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *store) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		list := make([]Todo, 0, len(s.todos))
		for _, t := range s.todos {
			list = append(list, t)
		}
		s.mu.Unlock()
		slices.SortFunc(list, func(a, b Todo) int { return a.ID - b.ID })
		writeJSON(w, http.StatusOK, list)
	})

	mux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		var t Todo
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil || t.Title == "" {
			writeError(w, http.StatusBadRequest, "title is required")
			return
		}
		s.mu.Lock()
		s.nextID++
		t.ID = s.nextID
		s.todos[t.ID] = t
		s.mu.Unlock()
		w.Header().Set("Location", fmt.Sprintf("/todos/%d", t.ID))
		writeJSON(w, http.StatusCreated, t)
	})

	mux.HandleFunc("GET /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			writeError(w, http.StatusBadRequest, "id must be a number")
			return
		}
		s.mu.Lock()
		t, ok := s.todos[id]
		s.mu.Unlock()
		if !ok {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		writeJSON(w, http.StatusOK, t)
	})

	mux.HandleFunc("PUT /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		var in Todo
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, ok := s.todos[id]; !ok {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		in.ID = id
		s.todos[id] = in
		writeJSON(w, http.StatusOK, in)
	})

	mux.HandleFunc("DELETE /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		s.mu.Lock()
		_, ok := s.todos[id]
		delete(s.todos, id)
		s.mu.Unlock()
		if !ok {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return logRequests(mux)
}

// statusRecorder remembers the status code a handler wrote.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// logRequests is middleware: it wraps every handler.
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		fmt.Printf("    [server] %-6s %-10s -> %d (%v)\n", r.Method, r.URL.Path, rec.status, time.Since(start).Round(time.Microsecond))
	})
}

// call is the client side: send a request, return status and decoded body.
func call(client *http.Client, method, url string, body any) (int, string) {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, url, &buf)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()
	var out bytes.Buffer
	out.ReadFrom(resp.Body)
	return resp.StatusCode, string(bytes.TrimSpace(out.Bytes()))
}

func main() {
	fmt.Println("Topic: REST API with net/http")

	ln, err := net.Listen("tcp", "127.0.0.1:0") // port 0: the OS picks a free port
	if err != nil {
		fmt.Fprintln(os.Stderr, "listen:", err)
		os.Exit(1)
	}
	srv := &http.Server{Handler: (&store{todos: map[int]Todo{}}).routes(), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintln(os.Stderr, "serve:", err)
		}
	}()
	base := "http://" + ln.Addr().String()
	fmt.Println("server listening on", base)

	client := &http.Client{Timeout: 3 * time.Second}
	type step struct {
		method, path string
		body         any
		want         int
	}
	failed := false
	for _, s := range []step{
		{"POST", "/todos", Todo{Title: "learn net/http"}, 201},
		{"POST", "/todos", Todo{Title: "write tests"}, 201},
		{"POST", "/todos", map[string]string{}, 400},
		{"GET", "/todos", nil, 200},
		{"GET", "/todos/1", nil, 200},
		{"PUT", "/todos/1", Todo{Title: "learn net/http", Done: true}, 200},
		{"GET", "/todos/abc", nil, 400},
		{"DELETE", "/todos/2", nil, 204},
		{"DELETE", "/todos/2", nil, 404},
		{"PATCH", "/todos/1", nil, 405},
		{"GET", "/todos", nil, 200},
	} {
		status, body := call(client, s.method, base+s.path, s.body)
		mark := "ok"
		if status != s.want {
			mark = fmt.Sprintf("MISMATCH (want %d)", s.want)
			failed = true
		}
		fmt.Printf("  %-6s %-10s %d %s  %s\n", s.method, s.path, status, mark, body)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "shutdown:", err)
		failed = true
	}
	fmt.Println("server shut down gracefully")

	if failed {
		os.Exit(1)
	}
	fmt.Println("all checks passed")
}
