# Webhook Service

Production-ready Webhook-as-a-Service platform built with Go and Svelte.

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Frontend   │────▶│  Backend    │────▶│ PostgreSQL  │
│ Svelte 5    │     │  Go/Fiber   │     │             │
│ Tailwind    │     │             │     └─────────────┘
└─────────────┘     │             │
                    │             │────▶ Redis (cache)
                    │             │
                    │             │────▶ NATS (queue)
                    └─────────────┘
                            │
                    ┌───────▼───────┐
                    │  Worker       │
                    │  (webhook     │
                    │   delivery)   │
                    └───────────────┘
```

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Backend | Go 1.22 + Fiber v2 |
| Frontend | Svelte 5 + SvelteKit + Tailwind CSS 4 |
| Database | PostgreSQL 17 |
| Cache | Redis 7 |
| Queue | NATS JetStream |
| Auth | JWT + API Keys |
| Deploy | Docker Compose |

## Quick Start

```bash
# Clone and start everything
cd webhook-service
docker-compose up --build

# Access:
# Frontend: http://localhost:5173
# API:      http://localhost:3000
# NATS UI:  http://localhost:8222
```

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/refresh` | Refresh token |

### Organizations
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/organizations` | List organizations |
| POST | `/api/v1/organizations` | Create organization |
| GET | `/api/v1/organizations/:id` | Get organization |
| PUT | `/api/v1/organizations/:id` | Update organization |
| DELETE | `/api/v1/organizations/:id` | Delete organization |
| GET | `/api/v1/organizations/:id/dashboard` | Dashboard stats |

### Applications
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/organizations/:id/applications` | List applications |
| POST | `/api/v1/organizations/:id/applications` | Create application |
| GET | `/api/v1/applications/:id` | Get application |
| PUT | `/api/v1/applications/:id` | Update application |
| DELETE | `/api/v1/applications/:id` | Delete application |

### Event Types
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/applications/:id/event-types` | List event types |
| POST | `/api/v1/applications/:id/event-types` | Create event type |
| DELETE | `/api/v1/applications/:id/event-types/:id` | Delete event type |

### Subscriptions
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/applications/:id/subscriptions` | List subscriptions |
| POST | `/api/v1/applications/:id/subscriptions` | Create subscription |
| GET | `/api/v1/subscriptions/:id` | Get subscription |
| PUT | `/api/v1/subscriptions/:id` | Update subscription |
| DELETE | `/api/v1/subscriptions/:id` | Delete subscription |
| POST | `/api/v1/subscriptions/:id/test` | Test subscription |

### Events
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/applications/:id/events` | Send event |
| GET | `/api/v1/applications/:id/events` | List events |
| GET | `/api/v1/events/:id` | Get event detail |

### Deliveries
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/applications/:id/deliveries` | List deliveries |
| GET | `/api/v1/deliveries/:id` | Get delivery detail |
| POST | `/api/v1/deliveries/:id/retry` | Retry delivery |

### Secrets
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/applications/:id/secrets` | List secrets |
| POST | `/api/v1/applications/:id/secrets` | Create secret |
| DELETE | `/api/v1/applications/:id/secrets/:id` | Delete secret |

## Authentication

Two methods supported:

### JWT Token
```bash
# Login first
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password"}'

# Use token in requests
curl http://localhost:3000/api/v1/organizations \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### API Key
```bash
# Create a secret first, then use it
curl http://localhost:3000/api/v1/applications/APP_ID/events \
  -H "X-API-Key: whkey_YOUR_API_KEY"
```

## Webhook Delivery

- **Retry strategy**: 3 attempts with exponential backoff (1s, 30s, 5min)
- **Signing**: HMAC-SHA256 signature in `X-Webhook-Signature` header
- **Dead Letter Queue**: Failed deliveries after max attempts go to DLQ

### Webhook Payload Format
```json
{
  "event_id": "uuid",
  "event_type": "order.created",
  "payload": { ... },
  "created_at": "2026-01-01T00:00:00Z"
}
```

### Webhook Headers
| Header | Description |
|--------|-------------|
| Content-Type | application/json |
| User-Agent | WebhookService/1.0 |
| X-Webhook-Event | Event type name |
| X-Webhook-Event-ID | Event UUID |
| X-Webhook-Signature | HMAC-SHA256 signature |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| SERVER_PORT | 3000 | API server port |
| DATABASE_URL | postgres://... | PostgreSQL connection |
| NATS_URL | nats://localhost:4222 | NATS connection |
| REDIS_URL | redis://localhost:6379 | Redis connection |
| JWT_SECRET | (change me) | JWT signing secret |
| JWT_EXPIRY | 24h | Token expiry |
| WORKER_COUNT | 4 | Delivery workers |
| FRONTEND_URL | http://localhost:5173 | CORS origin |

## Development

```bash
# Backend
cd backend
go mod download
go run ./cmd/server

# Frontend
cd frontend
npm install
npm run dev
```

## License

MIT
