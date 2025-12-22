## 1. Problem

**System**

- Single chat room
- One Go server
- React client

**Rules**

- User must provide a name before joining
- Name is immutable for the session
- Messages are in-memory only
- No persistence, no auth, no scaling
- If server restarts, everything is lost

This eliminates 90% of accidental complexity.

---

## 2. Minimal technical design (only what matters)

### Architecture

- Go HTTP server
- `/ws` WebSocket endpoint
- One global chat room (implicit)

### Server responsibilities

- Accept WebSocket connection **only if name is provided**
- Track connected clients `{name, conn}`
- Broadcast messages to all connected clients
- Remove client on disconnect

### Client responsibilities

- Block chat UI until name is entered
- Open WebSocket with name
- Send messages
- Render incoming messages

### Invariants (these become tests)

- No anonymous connections
- No empty messages
- Every message has a sender name
- Disconnected clients receive nothing

Stop here. No diagrams needed.

---

## 3. Translate design → backlog (behavioral, testable)

Each backlog item is a **promise to the user/system**.

### Backlog

1. User must provide a name before joining the chat
2. Server rejects WebSocket connections without a name
3. Connected users receive messages sent by other users
4. Messages include sender name and content
5. Empty messages are not broadcast
6. Disconnected users stop receiving messages

Notice:

- No Go structs
- No React state
- No channels
  This is intentional.

---

## 4. Translate backlog → tests (integration / boundary)

These are the tests you would write **before coding**.

### Test 1 — reject anonymous connection

**Given**

- Client opens WebSocket without name

**Expect**

- Connection is rejected (close or error)

---

### Test 2 — accept named connection

**Given**

- Client opens WebSocket with name `"Alice"`

**Expect**

- Connection succeeds
- Server associates client with `"Alice"`

---

### Test 3 — broadcast message

**Given**

- Alice and Bob are connected
- Alice sends `"Hello"`

**Expect**

- Bob receives `{ sender: "Alice", content: "Hello" }`

(Decide explicitly whether Alice also receives it; either is fine, but must be consistent.)

---

### Test 4 — reject empty message

**Given**

- Alice sends `""`

**Expect**

- No client receives a message

---

### Test 5 — disconnect cleanup

**Given**

- Alice disconnects
- Bob sends `"Hi"`

**Expect**

- Alice receives nothing
- Server remains stable

These tests define **truth**.
Any implementation that passes them is acceptable.

---

## 5. Only now: coding strategy (brief, non-prescriptive)

### Go

- One WebSocket handler
- One hub loop managing:

  - register
  - unregister
  - broadcast

- In-memory client map

### React

- Name gate (before socket connect)
- WebSocket lifecycle
- Append-only message list

You are free to:

- Refactor
- Change data structures
- Rewrite frontend logic

As long as tests pass, behavior is preserved.

---

## 6. What you intentionally do NOT test

- Internal goroutines
- Channel wiring
- React component internals
- CSS / layout

Those are implementation details.

---

## Final conclusion

With your constraints, the correct flow is:

> **Constraints → behavioral design → backlog promises → integration tests → code**

For this chat app:

- TDD ≈ **integration-first**
- Unit tests are optional and limited
- Complexity stays controlled because behavior is locked early

If you want next steps:

- Exact Go WebSocket integration test example
- How to simulate multiple clients in tests
- How much of this to skip if you’re just learning vs shipping
