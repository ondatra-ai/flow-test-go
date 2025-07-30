#!/bin/bash

# Helper script to run Go scripts in the scripts directory
# Usage: ./scripts/run-script.sh <script-name> [args...]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ $# -eq 0 ]; then
    echo "Usage: $0 <script-name> [args...]"
    echo ""
    echo "Available scripts:"
    echo "  get-pr-number                    - Get PR number for current branch"
    echo "  list-pr-conversations <pr-num>  - List unresolved PR conversations"
    echo "  resolve-pr-conversation <id>    - Resolve a PR conversation"
    echo ""
    echo "Examples:"
    echo "  $0 get-pr-number"
    echo "  $0 list-pr-conversations 123"
    echo "  $0 resolve-pr-conversation <thread-id> 'Fixed the issue'"
    echo ""
    echo "Note: You can also run scripts directly:"
    echo "  ./scripts/get-pr-number.go"
    echo "  ./scripts/list-pr-conversations.go 123"
    echo "  ./scripts/resolve-pr-conversation.go <thread-id>"
    exit 1
fi

SCRIPT_NAME="$1"
shift

case "$SCRIPT_NAME" in
    "get-pr-number")
        "$SCRIPT_DIR/get-pr-number.go" "$@"
        ;;
    "list-pr-conversations")
        "$SCRIPT_DIR/list-pr-conversations.go" "$@"
        ;;
    "resolve-pr-conversation")
        "$SCRIPT_DIR/resolve-pr-conversation.go" "$@"
        ;;
    *)
        echo "Unknown script: $SCRIPT_NAME"
        echo "Run '$0' without arguments to see available scripts."
        exit 1
        ;;
esac