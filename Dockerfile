# ---------- Stage 1: build React ----------
FROM node:24-alpine AS client-build
WORKDIR /app/client

COPY client/react/package.json client/react/pnpm-lock.yaml ./
RUN corepack enable && pnpm install

COPY client/react .
RUN pnpm build


# ---------- Stage 2: build Go ----------
FROM golang:1.25.5-alpine AS server-build
WORKDIR /app

COPY server/improved/go.mod server/improved/go.sum ./
RUN go mod download

COPY server/improved .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server


# ---------- Stage 3: runtime ----------
FROM alpine:latest
WORKDIR /app

# Go binary
COPY --from=server-build /app/server ./server

# React static files
COPY --from=client-build /app/client/dist ./static

EXPOSE 3001
CMD ["./server"]