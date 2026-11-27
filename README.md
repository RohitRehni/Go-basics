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
├── go.work                      # Workspace: 01-basics, 02-intermediate, 03-advanced
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
└── 03-advanced/
    ├── goroutines/waitgroup/, mutex/
    ├── channels/buffered-unbuffered/, select/, worker-pool/
    ├── context/timeout/, cancel/
    ├── sync-patterns/once/, atomic/
    ├── io/readers-writers/, json-encoding/
    └── reflection/reflect-basics/
```

Stages `04-dsa` and `05-integrations` are planned but not created yet.

## How to Run

### Run a lesson

Run the lesson folder (not the file) from the repo root:

```bash
go run ./01-basics/operators/arithmetic-op
go run ./02-intermediate/packages/bank-account
go run ./03-advanced/channels/worker-pool
```

Every lesson header has its exact `Run:` line. The concurrency lessons also
show a `Race:` line (`go run -race ...`) to check for data races.

### Run the tests

```bash
go test ./02-intermediate/testing/table-tests
go test -v ./02-intermediate/testing/table-tests   # show each case
```

### Check everything

```bash
for m in 01-basics 02-intermediate 03-advanced; do (cd $m && go vet ./...); done
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
./new-lesson.sh 04-dsa linked-lists singly-linked
```

## Learning Path

1. **01-basics**: Keywords, Operators, Arrays, Slices, Strings, Control Flow.
2. **02-intermediate**: Structs and methods, Interfaces, Error handling, Maps, Generics, Closures and defer, Packages, Testing.
3. **03-advanced**: Goroutines, Channels, Context, sync patterns, io and JSON, Reflection.
4. **04-dsa**: Standard algorithms and data structures in Go. (planned)
5. **05-integrations**: Connecting to external systems (SQL, Redis, HTTP). (planned)

## About setup-go-learning.sh

`setup-go-learning.sh` is the original generator. It is outdated: it puts
`main.go` straight into each concept folder (one program per concept) and it
deletes `go.work` and `README.md` when run again. Use `new-lesson.sh` instead.

## Requirements

- Go 1.26+ (`go.work` and each `go.mod` declare `go 1.26.4`; the lessons themselves need 1.25+ for `sync.WaitGroup.Go`)
- Built and tested with Go 1.26
