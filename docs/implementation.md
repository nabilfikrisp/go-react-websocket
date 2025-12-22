## Purpose

This document describes the **current implementation approach**.

---

## Architecture overview

- Go HTTP server
- Single `/ws` WebSocket endpoint
- One implicit global chat room

---

## Server implementation

- One WebSocket handler
- One hub loop responsible for:

  - Registering clients
  - Unregistering clients
  - Broadcasting messages

- In-memory map of connected clients `{name → connection}`
- No persistence layer

---

## Client implementation

### Lifecycle

1. Name gate renders first
2. After join:

   - WebSocket connection is created
   - Chat UI mounts

3. On unmount:

   - Socket is closed
   - No further state updates occur

---

### File responsibilities

| File                 | Responsibility                            |
| -------------------- | ----------------------------------------- |
| `name-gate.tsx`      | Name entry and join control               |
| `use-chat-socket.ts` | WebSocket lifecycle, message accumulation |
| `chat-room.tsx`      | Message rendering and user input          |

---

## React state strategy

- Message list is append-only
- No optimistic UI updates
- UI is a pure projection of hook state

---

## Notes

- Refactors are allowed
- Data structures may change
- Files may be split or merged

As long as the contract holds, the implementation is valid.

---

If you want next:

- a **diff checklist** for deciding SemVer PATCH vs MINOR
- or a **future `0.2.0` contract example** (e.g. persistence or rooms)
