# XM Companies Microservice

A microservice for managing company data with event-driven architecture using Clean Architecture principles.

## Tech Stack

- Go 1.24
- PostgreSQL 16
- Redpanda (Kafka-compatible event streaming)
- Docker & Docker Compose
- Make

## Architecture

- **Clean Architecture** - domain-centric design with clear layer separation
- **Outbox Pattern** - reliable event publishing with at-least-once delivery guarantee
- **JWT Authentication** - token-based API security with password validation

## Quick Start

1. Start all services:
```bash
make run
```

2. Verify health:
```bash
curl http://localhost:8080/health
```

3. Test the API:
```bash
# Generate token (authenticate with demo password)
export TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test","password":"demo-password-123"}' | jq -r '.token')

# Create company
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "XM Group",
    "description": "Financial services",
    "employees_count": 500,
    "registered": true,
    "type": "Corporations"
  }' | jq

# Get company
curl http://localhost:8080/api/v1/companies/{id} | jq
```

4. Verify events in Kafka (wait 5-6 seconds for outbox processor):
```bash
make kafka-consume
```

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/health` | No | Health check |
| POST | `/api/v1/auth/token` | Password | Generate JWT token (password: `demo-password-123`) |
| GET | `/api/v1/companies/{id}` | No | Get company |
| POST | `/api/v1/companies` | JWT | Create company |
| PATCH | `/api/v1/companies/{id}` | JWT | Update company |
| DELETE | `/api/v1/companies/{id}` | JWT | Delete company |

Full API specification: [api/openapi.yaml](api/openapi.yaml)

Request examples: [requests.md](requests.md)

## Authentication

The service uses a two-step authentication approach:

1. **Token Generation** (`/api/v1/auth/token`): Requires a password (`demo-password-123` for demo purposes)
   - In production, this would validate against a user database with hashed passwords
   - Returns a JWT token valid for 1 hour

2. **API Authorization**: Protected endpoints require the JWT token in the `Authorization` header
   - Format: `Authorization: Bearer <token>`
   - Token is validated on each request using HMAC signature

## Running Tests

```bash
make test
```

## How It Works

1. HTTP request → Handler validates input
2. Service executes business logic in transaction
3. Company + Event written to database (outbox table)
4. Transaction commits (both or neither)
5. Outbox Processor polls every 5s
6. Events published to Kafka
7. Outbox records marked as processed

**Result:** At-least-once delivery guarantee, no lost events even if Kafka is down.

**Event Ordering:** Hash balancing ensures events for same company_id go to same partition.

## Useful Commands

```bash
# Start services
make run

# Stop services
make down

# Restart app
make restart

# View logs
make logs

# Run tests
make test

# Run linter
make lint

# Check Kafka events
make kafka-consume

# Check companies in database
make db-companies

# Check outbox table
make db-outbox
```
