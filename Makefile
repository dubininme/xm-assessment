.PHONY: up down migrate-up migrate-down migrate-version generate lint lint-fix test test-unit test-integration logs kafka-consume db-companies db-outbox

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
	go test -tags=unit,integration ./... -race -v

test-unit:
	go test -tags=unit ./... -race -v

test-integration:
	go test -tags=integration ./... -race -v

kafka-consume:
	docker exec companies-kafka rpk topic consume company-events --num 10 --format json | jq

db-companies:
	docker exec companies-db psql -U xm_user -d xm_db -c "SELECT id, name, employees_count, registered, type FROM companies LIMIT 10;"

db-outbox:
	docker exec companies-db psql -U xm_user -d xm_db -c "SELECT id, event_type, aggregate_id, is_processed, created_at, processed_at FROM outbox ORDER BY id DESC LIMIT 10;"