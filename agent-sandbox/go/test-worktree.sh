#!/usr/bin/env bash
# Runs the worktree lifecycle with claude — create, list, delete — and prints
# what each command does. Nothing is asserted: read the output.
#
# Usage: ./test-worktree.sh

set -uo pipefail

cd -- "$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"

sandbox=(go run ./cmd/agent-sandbox)

run() {
  printf '\n$ agent-sandbox %s\n' "$*"
  "${sandbox[@]}" "$@" </dev/null
  printf '[exit %d]\n' "$?"
}

run tmp-run-test --agent claude --model opus "create a hello.txt file"

run worktree-list

run worktree-delete -b tmp-run-test --force

run worktree-list
