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

## 4. Translate backlog → tests (integration / boundary) for server

These are the tests you would write **before coding**.

### Test 1 — reject anonymous connection (done)

**Given**

- Client opens WebSocket without name

**Expect**

- Connection is rejected (close or error)

---

### Test 2 — accept named connection (done)

**Given**

- Client opens WebSocket with name `"Alice"`

**Expect**

- Connection succeeds
- Server associates client with `"Alice"`

---

### Test 3 — broadcast message (done)

**Given**

- Alice and Bob are connected
- Alice sends `"Hello"`

**Expect**

- Bob receives `{ sender: "Alice", content: "Hello" }`

(Alice also receives it)

---

### Test 4 — reject empty message (done)

**Given**

- Alice sends `""`

**Expect**

- No client receives a message

---

### Test 5 — disconnect cleanup (done)

**Given**

- Alice disconnects
- Bob sends `"Hi"`

**Expect**

- Alice receives nothing
- Server remains stable

These tests define **truth**.
Any implementation that passes them is acceptable.

---

## 4b. Translate backlog → executable behaviors (React client)

These are **runtime-verifiable behaviors**.
Each item must be _observable by running the app_ and _provably true_.

Treat each item as a checkbox you can execute, not speculate.

This is only executables, not tests.

---

## Executable 1 — name gate blocks chat

**Action**

- Open the app

**Observable result**

- Chat UI is not rendered
- Name input is visible
- **`useChatSocket` is not mounted**
- No WebSocket connection is opened

**Failure signal**

- Chat appears without a name
- Network tab shows a WebSocket connection
- Socket lifecycle logs appear before join

---

## Executable 2 — joining with name opens connection

**Action**

- Enter name `"Alice"`
- Click “Join”

**Observable result**

- Name input becomes disabled or disappears
- `useChatSocket` mounts
- WebSocket opens with `?name=Alice`
- Chat UI is rendered

**Failure signal**

- Name is editable after join
- Socket opens without name
- No socket connection

---

## Executable 3 — empty messages are not sent

**Action**

- Leave message input empty
- Press “Send”

**Observable result**

- `send()` is not invoked
- No WebSocket frame is sent
- `messages` array does not change

**Failure signal**

- Empty payload appears in Network tab
- Empty message appears in UI

---

## Executable 4 — sending a message

**Action**

- Type `"Hello"`
- Press “Send”

**Observable result**

- WebSocket sends `{ content: "Hello" }`
- Input clears immediately
- Message is **not optimistically added** by the UI

**Failure signal**

- Wrong payload shape
- Input not cleared
- UI mutates messages directly

---

## Executable 5 — receiving a message

**Action**

- Another client sends a message

**Observable result**

- `useChatSocket` appends message to `messages`
- Message appears in list
- Sender name and content are visible
- Order matches arrival order

**Failure signal**

- Message missing
- Sender not shown
- Order incorrect
- Chat UI mutates message history

---

## Executable 6 — disconnect handling

**Action**

- Stop the server or close the socket

**Observable result**

- `status` transitions to `closed` or `error`
- UI remains stable
- Message input and send button are disabled or status is shown

**Failure signal**

- App crashes
- User can invoke `send()` on a closed socket
- Messages are appended after disconnect

---

## Executable 7 — cleanup on unmount

**Action**

- Refresh the page or navigate away

**Observable result**

- WebSocket connection closes exactly once
- No lingering network activity
- No state updates after unmount

**Failure signal**

- Socket remains open
- Duplicate close attempts
- Console warnings about updates on unmounted components

---

## How to work with this

For each executable:

1. Implement minimal code
2. Run the app
3. Observe result
4. Check it off
5. Do not refactor yet

---

## Mapping executables → files (final)

**Message ownership rule (locked)**

- `useChatSocket` **owns and accumulates messages**
- `chat-room.tsx` is **purely presentational**
- No other file mutates message history

---

| Executable | File                 |
| ---------- | -------------------- |
| 1, 2       | `name-gate.tsx`      |
| 2, 6, 7    | `use-chat-socket.ts` |
| 3, 4, 5    | `chat-room.tsx`      |

---

## File contracts (now fully determined)

### `name-gate.tsx`

- Renders name input
- Emits `onJoin(name)`
- Never mounts or references socket logic
- Guarantees: no WebSocket exists while mounted

---

### `use-chat-socket.ts`

- Creates exactly one WebSocket per mount
- Accumulates messages in arrival order
- Exposes read-only `messages`
- Manages `status` (`connecting | open | closed | error`)
- Handles disconnect + cleanup

**Single source of truth for chat history**

---

### `chat-room.tsx`

- Receives `messages`, `status`, `send`
- Renders message list
- Guards empty sends
- Clears input on successful send
- Never mutates messages directly

**UI is a pure projection of hook state**

---

## Architectural invariants (must hold forever)

1. Messages only grow in `use-chat-socket`
2. No socket creation outside the hook
3. No UI component depends on socket internals
4. Unmount ⇒ socket closed exactly once

If any invariant breaks, at least one executable will fail.

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
