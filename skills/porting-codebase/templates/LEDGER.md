# LEDGER.md — {{SOURCE_LANG}} → {{TARGET_LANG}} port history

The port's **append-only** history. One line per event, newest last. `TASKS.md` §1 holds
the *current* status of each task; this file holds every transition that produced it.

**This file is written at the end of every task, in the same session that did the work.**
A task is not finished until its row in `TASKS.md` says so *and* its entries are here.

- Append only. Never edit, reorder, or delete an entry — not to tidy it, not to correct a
  typo in an older line, not to collapse a task's events into its latest state.
- `port`, `review`, `fix`, and `test` stay **separate entries**. That separation is what
  makes review yield measurable. Never fold them into one line.
- Every status transition is an entry: `TODO` on creation, `DOING` on claim, then `DONE`
  or `BLOCKED`. A reopened blocked task gets a new `TODO` and a new `DOING`; the earlier
  `BLOCKED` stays.
- Never reuse a task ID.
- Starting this file mid-port with no recorded history? Do not fabricate the missing
  events. Note the gap below and record every transition from that session onward.

Entry types: `TODO` · `DOING` · `BLOCKED` · `DONE` · `port` · `review` · `fix` · `test` ·
`prep` (scaffolding) · `verify` (an empirical check of a `PORTING.md` assumption) · `doc` ·
`amend` (a `PORTING.md` rule changed, with the rule ID and why).

---

## Entries

```
[T12] TODO      src/models/notification.py ← backend/models/notificacion.go
[T12] DOING     src/models/notification.py ← backend/models/notificacion.go
[T12] port      src/models/notification.py ← backend/models/notificacion.go (122 LOC, reads + writes)
[T12] review    src/models/notification.py — 3 findings: R22 log-and-continue inverted in the
                recipients loop; R19 NULL email surfaced as None instead of ""; R34 whitespace
                changed inside the UPDATE
[T12] fix       findings 1-2 applied; 3 rejected — the SQL text is byte-identical, only the
                Python string-literal indentation differs (models/notificacion.go:88)
[T12] test      tests/models/test_notification.py 9/9 · full suite 61/61
[T12] DONE      no new stubs · FINDINGS: +1 source defect (F7)
```

<!-- Replace the block above with this port's real entries. Keep them in one fenced block,
     or one block per wave — but never split a single task's entries across sections. -->

---

## History gaps

Only for ports that adopted this ledger after work had already started. Say what is
missing rather than inventing it.

| Task IDs | What is unrecorded | Why |
|---|---|---|
| | | |
