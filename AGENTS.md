# CLAUDE.md — server (SQC)

## Overview

Server Query Client (SQC) — a Go sidecar that queries game server stats via UDP extinfo and forwards them to chungusdb via HTTP.

- **Language**: Go (Gin)
- **Port**: 8080
- **Status**: Planned for deprecation (replaced by Lua→ENet→chungusway→gRPC pipeline)

## Endpoints

| Endpoint | Method | Auth | Purpose |
|----------|--------|------|---------|
| `/health` | GET | None | Health check |
| `/intermission` | GET | JWT Bearer | Triggers async stats export to chungusdb |

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `PLAYER_SERVICE_IP` | ChungusDB endpoint (e.g., `http://player:3000`) |
| `GAME_SERVER_IP` | Game server IP (usually `localhost` in shared namespace) |
| `GAME_SERVER_PORT` | UDP port for extinfo queries |
| `AUTH_SERVICE_IP` | Auth service endpoint |
| `SECRET_CHUNGUS` | JWT signing secret |
| `CHUNGUS_KEY` | API key for auth service |

## Key Files

| File | Purpose |
|------|---------|
| `main.go` | Entry point, router setup |
| `handler.go` | HTTP handlers |
| `middleware.go` | JWT validation middleware |
| `server_query.go` | Core logic: extinfo queries, JWT acquisition, stats export |
| `Dockerfile` | Multi-stage build (golang:1.24.5 → alpine) |

## Development

```bash
go build ./...
go run .
```

## Architecture Notes

- Runs as Docker sidecar alongside chungusmod (shared network namespace → reaches game server at `localhost`)
- Obtains JWT from auth at startup, caches it
- Stats exported to `POST /players/batch` on chungusdb (currently commented out on chungusdb side)
- Fire-and-forget: `/intermission` returns immediately, export runs in goroutine
- No tests, no graceful shutdown, debug prints throughout
