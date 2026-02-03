.PHONY: up down migrate-up migrate-down migrate-version generate lint lint-fix test logs kafka-consume db-companies db-outbox

OPENAPI_FILE = api/openapi.yaml
GEN_DIR = pkg/gen/oapi

run:
	docker-compose up -d --build

down:
	docker-compose down

restart:
	docker-compose restart app

logs:
	docker-compose logs -f app

migrate-up:
	docker-compose run --rm migrate up

migrate-down:
	docker-compose run --rm migrate down 1

migrate-down-all:
	docker-compose run --rm migrate down -all

migrate-version:
	docker-compose run --rm migrate version

generate:
	docker-compose exec app sh -c "mkdir -p pkg/gen/oapi && oapi-codegen -generate=types -package=oapi api/openapi.yaml > pkg/gen/oapi/types.go"

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test ./... -race -v

# Kafka: consume events from topic
kafka-consume:
	docker exec companies-kafka rpk topic consume company-events --num 10 --format json | jq

# Database: show companies
db-companies:
	docker exec companies-db psql -U xm_user -d xm_db -c "SELECT id, name, employees_count, registered, type FROM companies ORDER BY created_at DESC LIMIT 10;"

# Database: show outbox events
db-outbox:
	docker exec companies-db psql -U xm_user -d xm_db -c "SELECT id, event_type, aggregate_id, is_processed, created_at, processed_at FROM outbox ORDER BY id DESC LIMIT 10;"