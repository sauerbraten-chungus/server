# AGENTS.md — server (SQC)

## Overview

Server Query Client (SQC) — a Go sidecar that queries game server stats via UDP extinfo and forwards them to chungusdb via HTTP.

- **Language**: Go (Gin)
- **Port**: 8080
- **Status**: Planned for deprecation (replaced by Lua→ENet→chungusway→gRPC pipeline) — but still **unconditionally created per match** by chungustrator; deleting it today breaks match creation, not just stats

## Endpoints

| Endpoint | Method | Auth | Purpose |
|----------|--------|------|---------|
| `/health` | GET | None | Health check |
| `/intermission` | GET | JWT Bearer | Triggers async stats export to chungusdb |

## Environment Variables

All six are required — `main.go` reads them at startup and `log.Fatalf`s if `GAME_SERVER_PORT` is missing/non-numeric; the others silently break auth/stat calls when absent. See `.env.example` for a complete working dev set.

| Variable | Purpose |
|----------|---------|
| `PLAYER_SERVICE_IP` | ChungusDB endpoint (e.g., `http://host.docker.internal:3000`) |
| `GAME_SERVER_IP` | Game server IP (usually `localhost` in shared namespace) |
| `GAME_SERVER_PORT` | UDP port for extinfo queries |
| `AUTH_SERVICE_IP` | Auth service endpoint (e.g., `http://host.docker.internal:8081`) |
| `SECRET_CHUNGUS` | JWT signing secret |
| `CHUNGUS_KEY` | API key for auth service; must equal one of auth's `CHUNGUS_API_KEY_*` values |

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
