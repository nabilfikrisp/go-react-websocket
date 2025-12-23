## Purpose

This document describes the **current implementation approach**.

It explains how the system is structured today, without defining correctness rules.
As long as the contract holds, the implementation may change.

---

## Architecture overview

- Go HTTP server
- Single `/ws` WebSocket endpoint
- One implicit global chat room
- Event-based communication model

---

## Server implementation

- One WebSocket handler

- One hub loop responsible for:

  - Registering clients
  - Unregistering clients
  - Broadcasting events

- In-memory map of connected clients `{name → connection}`

- No persistence layer

### Event emission

- The server emits two categories of events:

  - User message events
  - System events (`join`, `leave`)

- All events are timestamped on the server before broadcast
- Join events are emitted after successful connection
- Leave events are emitted during disconnect cleanup

---

## Client implementation

### Lifecycle

1. Name gate renders first

2. After join:

   - WebSocket connection is created
   - Chat UI mounts
   - Event stream begins

3. On unmount:

   - Socket is closed
   - No further state updates occur

---

### File responsibilities

| File                 | Responsibility                             |
| -------------------- | ------------------------------------------ |
| `name-gate.tsx`      | Name entry and join control                |
| `use-chat-socket.ts` | WebSocket lifecycle and event accumulation |
| `chat-room.tsx`      | Event rendering and user input             |

---

## React state strategy

- Event list is append-only
- No optimistic UI updates
- UI is a pure projection of hook state
- Timestamps are treated as read-only data from the server

---

## Chat room projection

- Events are rendered in arrival order
- User message events are rendered as chat bubbles
- System events (join/leave) are visually distinct
- Own messages are visually distinct
- Consecutive messages from the same sender may be grouped

---

## Notes

- Refactors are allowed
- Data structures may change
- Files may be split or merged
- Internal naming and abstractions are flexible

As long as the contract holds, the implementation is valid.

---
