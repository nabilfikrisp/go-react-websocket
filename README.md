# Minimal WebSocket Chat

A deliberately minimal, single-room chat application built to **prioritize correctness and clarity over features**.

This repository is structured around a **behavioral contract**. Any implementation is valid as long as it satisfies the contract.

---

## What this is

- Single chat room
- One Go WebSocket server
- One React client
- No persistence, no auth, no scaling
- Server restart wipes all state

The goal is to eliminate accidental complexity and make system behavior explicit, testable, and observable.

---

## Repository structure

```

client/react # React client
server/improved # Go WebSocket server
docs/
contract.md # Behavioral contract (source of truth)
implementation.md # Current implementation details

```

---

## Source of truth

📌 **`docs/contract.md`**

This file defines:

- System rules and non-goals
- Behavioral guarantees (backlog)
- Server-side truth (integration / boundary tests)
- Client-side truth (runtime executables)
- Architectural invariants

The contract is **versioned by Git tags** using semantic versioning.

If the contract changes, the version changes.

---

## Implementation

📄 **`docs/implementation.md`**

Describes how the current code satisfies the contract:

- Architecture overview
- Server and client responsibilities
- File roles and structure

This document always reflects the **latest implementation** and may change freely as long as the contract holds.

---

## Running locally

### Server

```bash
cd server/improved
go run cmd/server/main.go
```

### Client

```bash
cd client/react
pnpm install
pnpm run dev
```

Open the client in your browser, enter a name, and start chatting.

---

## Versioning model

- `docs/contract@version.md`
- `implementation.md` always tracks HEAD

### SemVer rules

- **PATCH**: refactors, UI polish, internal changes
- **MINOR**: new observable behavior or guarantees
- **MAJOR**: breaking contract changes

---

## Design philosophy

- Behavior first
- Tests define truth
- UI is a projection of state
- No hidden side effects
- Refactor freely, break nothing observable

If an invariant breaks, at least one executable behavior must fail.

---

## Non-goals (by design)

- Multiple rooms
- Persistence
- Authentication
- Message history
- Production hardening

These are intentionally excluded.

---

## License

MIT (or your preferred license)
