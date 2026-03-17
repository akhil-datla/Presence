# Presence

**Simple, self-hosted attendance tracking API.**

[![CI](https://github.com/akhil-datla/Presence/actions/workflows/ci.yml/badge.svg)](https://github.com/akhil-datla/Presence/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/akhil-datla/Presence)](https://goreportcard.com/report/github.com/akhil-datla/Presence)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Presence is a lightweight, zero-dependency attendance tracking API. It runs as a single binary with an embedded SQLite database — no external services required.

## Features

- **Zero configuration** — runs out of the box with sensible defaults
- **Single binary** — no external database or services needed
- **JWT authentication** — secure user registration and login
- **Session management** — create and manage attendance sessions
- **Check-in/out tracking** — real-time attendance recording with timestamps
- **CSV export** — download attendance records as CSV
- **RESTful API** — clean, versioned JSON API
- **Docker ready** — multi-stage Dockerfile included
- **Graceful shutdown** — handles SIGINT/SIGTERM cleanly

## Quick Start

### From Source

```bash
git clone https://github.com/akhil-datla/Presence.git
cd Presence
make build
./presence
```

### With Docker

```bash
docker compose up -d
```

### Pre-built Binary

Download from [Releases](https://github.com/akhil-datla/Presence/releases) and run:

```bash
./presence
```

The server starts on `http://localhost:8080` by default.

## API Reference

All endpoints return JSON. Protected endpoints require `Authorization: Bearer <token>` header.

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register a new user |
| POST | `/api/v1/auth/login` | Login and get JWT token |

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/me` | Get current user profile |
| PUT | `/api/v1/users/me` | Update profile |
| DELETE | `/api/v1/users/me` | Delete account |

### Sessions

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/sessions` | Create a session |
| GET | `/api/v1/sessions` | List your sessions |
| GET | `/api/v1/sessions/:id` | Get session details |
| PUT | `/api/v1/sessions/:id` | Update a session |
| DELETE | `/api/v1/sessions/:id` | Delete a session |

### Attendance

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/sessions/:id/checkin` | Check in to session |
| POST | `/api/v1/sessions/:id/checkout` | Check out of session |
| GET | `/api/v1/sessions/:id/attendance` | List attendance records |
| GET | `/api/v1/sessions/:id/attendance/filter` | Filter by time |
| DELETE | `/api/v1/sessions/:id/attendance` | Clear attendance |
| GET | `/api/v1/sessions/:id/export/csv` | Export as CSV |

### Examples

**Register:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Jane","last_name":"Doe","email":"jane@example.com","password":"securepassword"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","password":"securepassword"}'
```

**Create Session:**
```bash
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"Team Standup"}'
```

**Check In:**
```bash
curl -X POST http://localhost:8080/api/v1/sessions/<session_id>/checkin \
  -H "Authorization: Bearer <token>"
```

## Configuration

| Variable | Flag | Default | Description |
|----------|------|---------|-------------|
| `PORT` | `--port` | `8080` | Server port |
| `DATABASE_PATH` | `--db` | `presence.db` | SQLite database path |
| `JWT_SECRET` | `--jwt-secret` | auto-generated | JWT signing secret |
| `LOG_LEVEL` | — | `info` | Log level |

Environment variables take precedence over flags.

## Development

```bash
make build    # Build binary
make run      # Build and run
make test     # Run tests
make lint     # Run linter
make docker   # Build Docker image
```

## License

MIT — see [LICENSE](LICENSE) for details.
