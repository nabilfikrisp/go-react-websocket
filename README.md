# Minimal WebSocket Chat

A deliberately minimal, single-room chat app.
Originally built while following an online Go course ([dasarpemrogramangolang.novalagung.com](https://dasarpemrogramangolang.novalagung.com/)), then improved bit by bit as a learning project.

## What this is

- One chat room
- Go WebSocket server
- React client
- No persistence, auth, or scaling
- Restart wipes state

## How it works

The system is defined by a behavioral contract.
As long as the contract holds, the implementation can change.

See `docs/contract.md` for the rules.

## Running locally

### With Docker (recommended)

```bash
docker compose up
```

This builds and runs both the Go WebSocket server and the React client.

### Without Docker (optional)

#### Server

```bash
cd server/improved
go run cmd/server/main.go
```

#### Client

```bash
cd client/react
pnpm install
pnpm run dev
```
