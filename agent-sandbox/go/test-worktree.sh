#!/usr/bin/env bash
# Runs the worktree lifecycle with claude on golang:1.26-alpine — create, list,
# delete — and shows that the state file, not git, decides what the verbs can
# see. It also proves that the external image's Go toolchain is available.
# Nothing is asserted: read the output.
#
# Needs Docker and an authenticated claude. Point XDG_CONFIG_HOME at a scratch
# directory to keep the real ~/.config/agent-sandbox out of it.
#
# Usage: ./test-worktree.sh

set -uo pipefail

cd -- "$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"

sandbox=(go run ./cmd/agent-sandbox)
state="${XDG_CONFIG_HOME:-$HOME/.config}/agent-sandbox/worktrees.jsonl"

run() {
  printf '\n$ agent-sandbox %s\n' "$*"
  "${sandbox[@]}" "$@" </dev/null
  printf '[exit %d]\n' "$?"
}

# git is run directly for the setup a sandbox run would never do.
plain_git() {
  printf '\n$ git %s\n' "$*"
  git "$@"
  printf '[exit %d]\n' "$?"
}

show_state() {
  printf '\n$ cat %s\n' "$state"
  cat -- "$state" 2>/dev/null || printf '(no state file)\n'
}

# One sandbox run: builds or reuses the image derived from golang:1.26-alpine,
# creates ../tmp-run-test, and appends its line to the state file. It must show
# Go 1.26 before creating hello.txt. A second run reuses the derived image.
run tmp-run-test --agent claude --model opus --base-image golang:1.26-alpine \
  "run go version, then create a file named hello.txt at the repository root containing the text hello world"
show_state

# A worktree the sandbox did not create. Git knows it, the state file does not,
# so nothing below may list, delete or even refuse to enter it.
plain_git worktree add -b tmp-hand-test ../tmp-hand-test

run worktree-list                     # tmp-run-test only
run worktree-delete -b tmp-hand-test  # refused: no sandbox worktree on that branch

# Without --yes the prompt reads EOF from </dev/null, which counts as no.
run worktree-delete-all
run worktree-delete-all --yes

run worktree-list                     # nothing left
show_state                            # and no line left either

plain_git worktree list               # tmp-hand-test survived it all

# Cleaning up the hand-made worktree by hand, which is now the only way.
plain_git worktree remove --force ../tmp-hand-test
plain_git branch -D tmp-hand-test
