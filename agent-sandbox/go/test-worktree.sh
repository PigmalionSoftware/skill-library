#!/usr/bin/env bash
# Runs every agent-sandbox worktree command and its normal execution flags.
# The default path uses claude on golang:1.26-alpine and covers create, resume,
# list, editor, single delete, and bulk delete. It shows that the state file,
# not git, decides what the verbs can see.
# Nothing is asserted: read the output.
#
# Needs Docker and an authenticated claude. Point XDG_CONFIG_HOME at a scratch
# directory to keep the real ~/.config/agent-sandbox out of it.
# TEST_EDITOR=1 opens VS Code for the resumed worktree. TEST_CLAUDE_IMAGES=1
# or TEST_CODEX_IMAGES=1 needs the corresponding authenticated agent and
# verifies actual --image attachment handling.
#
# Usage: ./test-worktree.sh

set -uo pipefail

cd -- "$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"

sandbox=(go run ./cmd/agent-sandbox)
state="${XDG_CONFIG_HOME:-$HOME/.config}/agent-sandbox/worktrees.jsonl"
agent="${AGENT:-claude}"
model="${MODEL:-opus}"
base_image="${BASE_IMAGE:-golang:1.26-alpine}"
prompt_file="$(mktemp)"
image_file="$(mktemp --suffix=.png)"

# A real 1x1 PNG exercises the agent's visual attachment path. An empty file
# would only prove that Docker mounted something, not that Claude or Codex can
# decode an image.
base64 --decode >"$image_file" <<'PNG'
iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAusB9Y9Jx60AAAAASUVORK5CYII=
PNG

cleanup() {
  rm -f -- "$prompt_file" "$image_file"
}
trap cleanup EXIT

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

run --version
run --help
run resume --help
run worktree-list --help
run worktree-editor --help
run worktree-delete --help
run worktree-delete-all --help

# The initial run builds or reuses the image derived from golang:1.26-alpine,
# creates tmp-run-test, and appends one line to the state file. It must show Go
# 1.26 before creating hello.txt.
run -b tmp-run-test -a claude -m opus -i golang:1.26-alpine \
  "run go version, then create a file named hello.txt at the repository root containing the text hello world"
show_state

# Resume identifies the same worktree by the branch name recorded in the state
# file. It must see the uncommitted hello.txt from the first session, add a
# second file beside it, and leave the state file with its original one line.
run resume -b tmp-run-test -a claude -m opus -i golang:1.26-alpine \
  "read hello.txt, then create resumed.txt at the repository root containing the text resumed successfully"
show_state

# --file-prompt takes the task from a regular host file. This run is kept for
# worktree-delete-all below.
printf '%s\n' "create file-prompt.txt at the repository root containing file prompt works" >"$prompt_file"
run -b tmp-file-prompt-test -a "$agent" -m "$model" -f "$prompt_file"

# worktree-delete handles one clean sandbox worktree directly. The agent is
# asked only to inspect it, so deletion should not need --force.
run -b tmp-delete-test -a "$agent" -m "$model" -i "$base_image" \
  "run git status and do not modify any repository file"
run worktree-delete -b tmp-delete-test

# Opening an editor is intentionally opt-in because it launches a host GUI.
if [[ "${TEST_EDITOR:-0}" == "1" ]]; then
  run worktree-editor -b tmp-run-test
else
  printf '\n[skipping worktree-editor; set TEST_EDITOR=1 to launch VS Code]\n'
fi

# Claude receives mounted images as visual paths in its initial prompt. This is
# opt-in because the normal worktree lifecycle needs only existing Claude auth.
if [[ "${TEST_CLAUDE_IMAGES:-0}" == "1" ]]; then
  run -b tmp-claude-image-test -a claude -m opus -i golang:1.26-alpine \
    --image "$image_file" -- \
    "inspect the attached image as visual input, then create claude-image-flag.txt at the repository root containing image flag works"
else
  printf '\n[skipping Claude --image test; set TEST_CLAUDE_IMAGES=1]\n'
fi

# Codex uses its native image arguments. It stays opt-in because it needs its
# own host authentication in addition to the default Claude lifecycle.
if [[ "${TEST_CODEX_IMAGES:-0}" == "1" ]]; then
  run -b tmp-image-test -a codex --image "$image_file" -- \
    "create image-flag.txt at the repository root containing image flag works"
else
  printf '\n[skipping Codex --image test; set TEST_CODEX_IMAGES=1]\n'
fi

# A worktree the sandbox did not create. Git knows it, the state file does not,
# so nothing below may list, delete or even refuse to enter it.
plain_git worktree add -b tmp-hand-test ../tmp-hand-test

run worktree-list                     # tmp-run-test only
run worktree-delete -b tmp-hand-test  # refused: no sandbox worktree on that branch

# Resume follows the same state-file boundary: Git knows tmp-hand-test, but
# the sandbox never recorded it, so this must fail before Docker starts.
run resume -b tmp-hand-test -a claude -m opus -i golang:1.26-alpine \
  "this must not start"

# Without --yes the prompt reads EOF from </dev/null, which counts as no.
run worktree-delete-all
run worktree-delete-all --yes

run worktree-list                     # nothing left
show_state                            # and no line left either

plain_git worktree list               # tmp-hand-test survived it all

# Cleaning up the hand-made worktree by hand, which is now the only way.
plain_git worktree remove --force ../tmp-hand-test
plain_git branch -D tmp-hand-test
