# Go Learning and Modules

A structured Go learning journey using **Go Modules** and a **Go Workspace**.
Rust-Basics and Python-Basics follow the same stage and concept layout.

## Project Structure

Each numbered folder is a learning stage and its own module (`learning/<stage>`),
listed in `go.work`. Inside a stage, every concept has a folder, and every
**lesson is its own folder with a `main.go`**. In Go a folder is a package, so
one folder per lesson is what lets each lesson have its own `main()` and still
pass `go vet` and editor checks.

```text
.
├── go.work                      # Workspace: all five stage modules
├── Go-basics.code-workspace     # VS Code workspace file
├── README.md                    # This file
├── new-lesson.sh                # Adds a lesson (and a stage, if new)
├── setup-go-learning.sh         # Original generator (outdated, see below)
├── 01-basics/
│   ├── arrays/array-op/
│   ├── control-flow/            # empty: yours to fill
│   ├── keywords/hello/
│   ├── operators/arithmetic-op/, comparison-op/
│   ├── slices/slice-op/
│   └── strings/palindrome/, roman-to-int/, longest-common-prefix/
├── 02-intermediate/
│   ├── structs/define-structs/, methods/, embedding/
│   ├── interfaces/interfaces/, type-switch/, stringer/
│   ├── error-handling/errors-basics/, wrapping/, custom-errors/
│   ├── maps/maps/
│   ├── generics/generic-funcs/, constraints/
│   ├── closures-defer/closures/, defer-panic-recover/
│   ├── packages/bank-account/   # multi-file: main.go, report.go, account/account.go
│   └── testing/table-tests/     # main.go + main_test.go
├── 03-advanced/
│   ├── goroutines/waitgroup/, mutex/
│   ├── channels/buffered-unbuffered/, select/, worker-pool/
│   ├── context/timeout/, cancel/
│   ├── sync-patterns/once/, atomic/
│   ├── io/readers-writers/, json-encoding/
│   └── reflection/reflect-basics/
├── 04-dsa/                      # standard library only; every lesson checks its own results
│   ├── searching/linear-search/, binary-search/
│   ├── sorting/simple-sorts/, merge-sort/, quick-sort/
│   ├── linked-lists/singly-linked-list/, doubly-linked-list/
│   ├── stacks-queues/stack/, queue-deque/
│   ├── hash-tables/hash-map/
│   ├── trees/binary-search-tree/, tree-traversals/
│   ├── heaps/min-heap/
│   ├── graphs/bfs-dfs/, dijkstra/, topological-sort/
│   ├── dynamic-programming/fibonacci-memo/, knapsack/, longest-common-subsequence/
│   ├── backtracking/permutations-subsets/, n-queens/
│   └── techniques/two-pointers/, sliding-window/
└── 05-integrations/             # uses pgx, go-redis and x/net/html
    ├── compose.yaml             # local Postgres + Redis for these lessons
    ├── database/postgres-crud/, transactions/
    ├── redis/redis-basics/, redis-cache/
    ├── http/http-client/, rest-api/
    ├── json/json-encoding/
    ├── scraping/scrape-pages/
    └── capstone/notes-api/      # multi-file: main.go, handlers.go, store.go, cache.go
```

## How to Run

### Run a lesson

Run the lesson folder (not the file) from the repo root:

```bash
go run ./01-basics/operators/arithmetic-op
go run ./02-intermediate/packages/bank-account
go run ./03-advanced/channels/worker-pool
go run ./04-dsa/graphs/dijkstra
go run ./05-integrations/capstone/notes-api      # needs the services below
```

The `04-dsa` lessons implement each structure from scratch, then compare the
result with the standard library (or a brute-force answer). Each one prints
`all checks passed` and exits 1 if anything is off.

Every lesson header has its exact `Run:` line. The concurrency lessons also
show a `Race:` line (`go run -race ...`) to check for data races.

### Run the tests

```bash
go test ./02-intermediate/testing/table-tests
go test -v ./02-intermediate/testing/table-tests   # show each case
```

### Check everything

```bash
for m in 01-basics 02-intermediate 03-advanced 04-dsa 05-integrations; do (cd $m && go vet ./...); done
gofmt -l .   # prints nothing when all files are formatted
```

### Add a new lesson

```bash
./new-lesson.sh 01-basics control-flow if-else
go run ./01-basics/control-flow/if-else
```

This creates `01-basics/control-flow/if-else/main.go` from a starter template.
Lesson names are kebab-case.

### Multi-file lessons

A lesson folder can hold more than one file:

- **More `.go` files in the same folder** share `package main`, so their
  functions are visible to `main.go` without an import (see `report.go` in
  `02-intermediate/packages/bank-account`).
- **A sub-folder is a separate package** you import by its path, for example
  `learning/02-intermediate/packages/bank-account/account`. Only Capitalised
  names are visible outside it. This is how Go does object-oriented code: a
  type, a `NewXxx` constructor and methods, with private fields kept inside
  the package.

### Add a new stage

Use the same command with a new stage folder. The script runs `go mod init
learning/<stage>` and `go work use` for you:

```bash
./new-lesson.sh 06-microservices grpc hello-grpc
```

## Services for 05-integrations

The database, Redis and capstone lessons need PostgreSQL and Redis. Start both
with Docker:

```bash
docker compose -f 05-integrations/compose.yaml up -d --wait   # start
docker compose -f 05-integrations/compose.yaml down           # stop
```

| Service | Address | Env var (default) |
|---|---|---|
| PostgreSQL 16 | `127.0.0.1:55432`, user/password/db `learning` | `DATABASE_URL=postgres://learning:learning@127.0.0.1:55432/learning` |
| Redis | `127.0.0.1:56379` | `REDIS_URL=redis://127.0.0.1:56379/0` |

Some things to know about the setup:

- **Separate ports:** it uses 55432 and 56379 instead of the usual 5432 and
  6379, so it never touches another Postgres or Redis already running on this
  machine.
- **Throwaway data:** data is kept in memory (tmpfs), so every restart starts
  empty. Each lesson also creates its own `go_*` tables and `go:*` Redis keys
  and removes them when it finishes, so lessons can be run again.
- **Shared stack:** the same `compose.yaml` (project `learning-services`) sits
  in Rust-Basics, Python-Basics and OOPerl-basics, so one running stack serves
  all four repos.
- **Errors:** if a service is not running, a lesson prints one line telling
  you to start it and exits 1.
- **Local only:** the HTTP, JSON and scraping lessons need no services and no
  internet. Each one starts its own server on `127.0.0.1`.

Dependencies of the `05-integrations` module, recorded in its `go.mod` and
downloaded on the first `go run`:

- `github.com/jackc/pgx/v5`: PostgreSQL driver and connection pool
- `github.com/redis/go-redis/v9`: Redis client
- `golang.org/x/net/html`: HTML parser for scraping

## Learning Path

1. **01-basics**: Keywords, Operators, Arrays, Slices, Strings, Control Flow.
2. **02-intermediate**: Structs and methods, Interfaces, Error handling, Maps, Generics, Closures and defer, Packages, Testing.
3. **03-advanced**: Goroutines, Channels, Context, sync patterns, io and JSON, Reflection.
4. **04-dsa**: Searching, Sorting, Linked lists, Stacks and queues, Hash tables, Trees, Heaps, Graphs, Dynamic programming, Backtracking, Two pointers and sliding window.
5. **05-integrations**: PostgreSQL (CRUD, transactions), Redis (basics, caching), HTTP client and REST API, JSON, Web scraping, and a capstone notes API that combines them.

## About setup-go-learning.sh

`setup-go-learning.sh` is the original generator. It is outdated: it puts
`main.go` straight into each concept folder (one program per concept) and it
deletes `go.work` and `README.md` when run again. Use `new-lesson.sh` instead.

## Requirements

- Go 1.26+ (`go.work` and each `go.mod` declare `go 1.26.4`; the lessons themselves need 1.25+ for `sync.WaitGroup.Go`)
- Built and tested with Go 1.26
- Docker with Compose, for the `05-integrations` database and Redis lessons
