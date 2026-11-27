package main

// handlers.go: the HTTP layer. Only this file knows about requests and status codes.

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type API struct {
	store *NoteStore
	cache *NoteCache
}

func NewAPI(store *NoteStore, cache *NoteCache) *API { return &API{store: store, cache: cache} }

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /notes", a.list)
	mux.HandleFunc("POST /notes", a.create)
	mux.HandleFunc("GET /notes/{id}", a.get)
	mux.HandleFunc("PUT /notes/{id}", a.update)
	mux.HandleFunc("DELETE /notes/{id}", a.delete)
	return mux
}

type noteInput struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// fail maps store errors to HTTP responses; unexpected errors are logged, not leaked.
func fail(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	fmt.Fprintln(os.Stderr, "internal error:", err)
	writeError(w, http.StatusInternalServerError, "internal error")
}

func pathID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "id must be a positive number")
		return 0, false
	}
	return id, true
}

func readInput(w http.ResponseWriter, r *http.Request) (noteInput, bool) {
	var in noteInput
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)) // cap the body at 1 MB
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil || strings.TrimSpace(in.Title) == "" {
		writeError(w, http.StatusBadRequest, "JSON with a non-empty title is required")
		return in, false
	}
	return in, true
}

func (a *API) list(w http.ResponseWriter, r *http.Request) {
	notes, err := a.store.List(r.Context())
	if err != nil {
		fail(w, err)
		return
	}
	writeJSON(w, http.StatusOK, notes)
}

func (a *API) create(w http.ResponseWriter, r *http.Request) {
	in, ok := readInput(w, r)
	if !ok {
		return
	}
	n, err := a.store.Create(r.Context(), in.Title, in.Body)
	if err != nil {
		fail(w, err)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("/notes/%d", n.ID))
	writeJSON(w, http.StatusCreated, n)
}

// get is the read-through path: Redis first, then Postgres (and fill Redis).
func (a *API) get(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if n, hit := a.cache.Get(r.Context(), id); hit {
		w.Header().Set("X-Cache", "HIT")
		writeJSON(w, http.StatusOK, n)
		return
	}
	n, err := a.store.Get(r.Context(), id)
	if err != nil {
		fail(w, err)
		return
	}
	a.cache.Set(r.Context(), n)
	w.Header().Set("X-Cache", "MISS")
	writeJSON(w, http.StatusOK, n)
}

func (a *API) update(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	in, ok := readInput(w, r)
	if !ok {
		return
	}
	n, err := a.store.Update(r.Context(), id, in.Title, in.Body)
	if err != nil {
		fail(w, err)
		return
	}
	a.cache.Invalidate(r.Context(), id) // database first, then drop the stale copy
	writeJSON(w, http.StatusOK, n)
}

func (a *API) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := a.store.Delete(r.Context(), id); err != nil {
		fail(w, err)
		return
	}
	a.cache.Invalidate(r.Context(), id)
	w.WriteHeader(http.StatusNoContent)
}
