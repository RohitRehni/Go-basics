#!/bin/bash

# ==============================================================================
# Add a Go lesson
# ==============================================================================
# Usage:
#   ./new-lesson.sh <stage> <concept> <lesson>
#   ./new-lesson.sh 01-basics control-flow if-else
#   ./new-lesson.sh 02-intermediate structs bank-account
#
# Every lesson is its own folder with a main.go (in Go a folder is a package,
# so one folder = one program). Add more .go files to the same folder for a
# multi-file lesson; they share package main. A new stage becomes its own
# module (learning/<stage>) and is added to go.work.
# Run a lesson with: go run ./<stage>/<concept>/<lesson>
# ==============================================================================

set -euo pipefail

if [ $# -ne 3 ]; then
    sed -n '6,9p' "$0" | sed 's/^# //'
    exit 1
fi

STAGE="$1"
CONCEPT="$2"
LESSON="$3"

cd "$(dirname "$0")"

if [[ ! $LESSON =~ ^[a-z][a-z0-9-]*$ ]]; then
    echo "Lesson name must be kebab-case, e.g. if-else" >&2
    exit 1
fi

# --- New stage: its own module, added to the workspace ---
if [ ! -f "$STAGE/go.mod" ]; then
    mkdir -p "$STAGE"
    (cd "$STAGE" && go mod init "learning/$STAGE" >/dev/null 2>&1)
    go work use "./$STAGE"
    echo "Created stage $STAGE (module learning/$STAGE)"
fi

DIR="$STAGE/$CONCEPT/$LESSON"
if [ -e "$DIR" ]; then
    echo "$DIR already exists" >&2
    exit 1
fi

mkdir -p "$DIR"
cat > "$DIR/main.go" <<EOF
// Topic: $LESSON | Concept: $CONCEPT | Stage: $STAGE
// Run:   go run ./$DIR
//
// Explore:
// * TODO: list what you want to try here
package main

import "fmt"

func main() {
	fmt.Println("Hello from $LESSON in $STAGE/$CONCEPT!")
	// TODO: add your learning code here
}
EOF

echo "Created $DIR/main.go"
echo "Run it: go run ./$DIR"
