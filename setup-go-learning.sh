#!/bin/bash

# ==============================================================================
# Go Learning Workspace Generator
# ==============================================================================
# Description:
#   Creates a modern Go Modules + Workspace structure for learning Go.
#   Organizes topics from Basics -> Advanced -> DSA -> Integrations.
#   Automatically generates documentation (README.md) within the project.
#
# Usage:
#   ./setup-go-learning.sh [project_name]
#
# Arguments:
#   project_name : Optional. Name of the root folder (Default: "go-learning-journey")
#
# Created: June 29, 2026
# ==============================================================================

set -e # Exit immediately if a command exits with a non-zero status

# --- Configuration ---
PROJECT_NAME="${1:-go-learning-journey}"
MODULE_PREFIX="learning"

# Define the learning modules and their sub-topics
# Format: "folder_name:sub_topic1,sub_topic2"
declare -a MODULES=(
    "01-basics:keywords,strings,control-flow"
    "02-intermediate:structs,interfaces,error-handling"
    "03-advanced:goroutines,channels,context"
    "04-dsa:arrays-slices,maps,linked-lists,trees"
    "05-integrations:database,redis,api-client"
)

echo "Starting Go Learning Workspace Setup..."
echo "Project Name: $PROJECT_NAME"

# 1. Create Root Directory
if [ -d "$PROJECT_NAME" ]; then
    echo " Directory '$PROJECT_NAME' already exists. Cleaning up old Go files..."
    rm -f "$PROJECT_NAME/go.work"
    rm -f "$PROJECT_NAME/README.md"
else
    mkdir -p "$PROJECT_NAME"
fi
cd "$PROJECT_NAME"

# 2. Initialize Go Workspace
echo "Initializing Go Workspace (go.work)..."
go work init

# 3. Create Modules and Sub-directories
echo "Creating Modules and Directory Structure..."

for module_entry in "${MODULES[@]}"; do
    # Split module name and sub-topics
    module_name="${module_entry%%:*}"
    sub_topics="${module_entry#*:}"
    
    # Create module directory
    mkdir -p "$module_name"
    cd "$module_name"
    
    # Initialize go.mod
    go mod init "$MODULE_PREFIX/$module_name"
    
    # Add module to workspace (from parent dir context, but we are inside, so go up)
    # Actually, 'go work use' works best from the root. We will do it at the end or use relative path.
    # Let's add it to the workspace from the root perspective.
    cd ..
    go work use "./$module_name"
    
    # Create sub-topic directories and placeholder main.go files
    IFS=',' read -ra TOPICS <<< "$sub_topics"
    for topic in "${TOPICS[@]}"; do
        topic_path="$module_name/$topic"
        mkdir -p "$topic_path"
        
        # Create a placeholder main.go
        cat > "$topic_path/main.go" <<EOF
package main

import "fmt"

// Topic: $topic
// Module: $module_name
func main() {
    fmt.Println("Hello from $topic in $module_name!")
    // TODO: Add your learning code here
}
EOF
    done
    
    echo "  Created module: $module_name"
done

# 4. Generate Documentation (README.md)
echo "📝 Generating Documentation (README.md)..."

cat > README.md <<EOF
# $PROJECT_NAME

A structured Go learning journey using **Go Modules** and **Workspaces**.

## Created
Generated on: $(date +"%Y-%m-%d")

## Project Structure

This project uses a **Go Workspace** (\`go.work\`) to manage multiple independent modules.
Each numbered folder represents a learning stage.

\`\`\`text
.
├── go.work              # Workspace definition
├── README.md            # This file
├── 01-basics/           # Syntax, Keywords, Strings
├── 02-intermediate/     # Structs, Interfaces, Errors
├── 03-advanced/         # Concurrency (Goroutines, Channels)
├── 04-dsa/              # Data Structures & Algorithms
└── 05-integrations/     # DB, Redis, API clients
\`\`\`

## How to Run

### Run a specific topic
You can run any \`main.go\` directly from the root using the workspace:

\`\`\`bash
# Example: Run the 'keywords' topic from basics
go run ./01-basics/keywords/main.go

# Example: Run the 'channels' topic from advanced
go run ./03-advanced/channels/main.go
\`\`\`

### Add a new topic
1. Create a folder inside the relevant module (e.g., \`01-basics/my-new-topic\`).
2. Create a \`main.go\` file inside it.
3. Run it: \`go run ./01-basics/my-new-topic/main.go\`.

### Add a new Module (Stage)
If you need a completely new stage (e.g., \`06-microservices\`):
\`\`\`bash
mkdir 06-microservices
cd 06-microservices
go mod init learning/06-microservices
cd ..
go work use ./06-microservices
\`\`\`

## Learning Path

1. **01-basics**: Variables, Data Types, Control Flow.
2. **02-intermediate**: Custom Types, Interfaces, Error Handling patterns.
3. **03-advanced**: Concurrency primitives, Context, Sync package.
4. **04-dsa**: Implementation of standard algorithms in Go.
5. **05-integrations**: Connecting to external systems (SQL, NoSQL, HTTP).

## 🛠️ Requirements
- Go 1.18+ (Workspace support)
- Latest Go version recommended (1.22+)

EOF

echo "Setup Complete!"
echo ""
echo "Next steps:"
echo "  cd $PROJECT_NAME"
echo "  go run ./01-basics/keywords/main.go"
echo ""
echo "Read the README.md for full documentation."   