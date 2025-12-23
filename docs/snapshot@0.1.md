## 1. Problem

### System

- Single chat room
- One Go server
- React client

### Rules

- User must provide a name before joining
- Name is immutable for the session
- Messages are in-memory only
- No persistence, no authentication, no scaling
- Server restart clears all state

**Non-goals**

- Multi-room support
- Message history persistence
- User identity beyond session name
- Reliability across restarts

---

## 2. Behavioral guarantees

1. A user must provide a name before joining the chat
2. The server rejects WebSocket connections without a name
3. Connected users receive messages sent by other users
4. Messages include sender name and content
5. Empty messages are not broadcast
6. Disconnected users stop receiving messages

All guarantees are externally observable.

---

## 3. Message model (wire contract)

```ts
type ChatMessage = {
  sender: string;
  content: string;
};
```

---

## 4. Server truth (integration / boundary tests)

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
- Alice also receives the message

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

Any implementation that passes these tests is compliant.

---

## 5. Client truth (runtime executables)

### Executable 1 — name gate blocks chat

- Chat UI is not rendered before name entry
- No WebSocket exists before join

---

### Executable 2 — joining with name opens connection

- WebSocket opens with `?name=Alice`
- Chat UI renders
- Name becomes immutable

---

### Executable 3 — empty messages are not sent

- `send()` is not invoked
- No WebSocket frame is sent
- Messages do not change

---

### Executable 4 — sending a message

- WebSocket sends `{ content: "Hello" }`
- Input clears immediately
- UI does not optimistically append messages

---

### Executable 5 — receiving a message

- Messages are appended in arrival order
- Sender name and content are visible

---

### Executable 6 — disconnect handling

- Status becomes `closed` or `error`
- UI remains stable
- Sending is disabled

---

### Executable 7 — cleanup on unmount

- WebSocket closes exactly once
- No state updates after unmount

---

## 6. Ownership and file contracts

### Message ownership (locked)

- `use-chat-socket` owns and accumulates message history
- No other module mutates messages

---

### `name-gate.tsx`

- Renders name input
- Emits `onJoin(name)`
- Does not create or reference sockets
- Guarantees no WebSocket exists while mounted

---

### `use-chat-socket.ts`

- Creates exactly one WebSocket per mount
- Accumulates messages in arrival order
- Exposes read-only `messages`
- Manages `status` (`connecting | open | closed | error`)
- Handles disconnect and cleanup

---

### `chat-room.tsx`

- Receives `messages`, `status`, `send`
- Guards empty sends
- Clears input on successful send
- Never mutates message history

---

## 7. Architectural invariants

1. Messages only grow inside `use-chat-socket`
2. No socket creation outside the hook
3. UI components do not depend on socket internals
4. Unmount closes the socket exactly once

If any invariant breaks, at least one executable fails.
