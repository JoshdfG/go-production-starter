-include .env
export

.PHONY: build run test test-v test-all lint fmt tidy swagger \
        docker-up docker-down migrate ps verify \
        register login create list fail

# ── build ────────────────────────────────────────────────
build:
	go build -v -o bin/app ./cmd/app

run:
	go run ./cmd/app

# ── test ─────────────────────────────────────────────────
test:
	go test ./internal/usecase/... -race -count=1

test-v:
	go test ./internal/usecase/... -v -race -count=1

test-all:
	go test ./... -race -count=1 -timeout 30s

test-integration:
	go test ./internal/repo/... -v -run TestPostgresSuite -timeout 60s

# ── quality ──────────────────────────────────────────────
lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

tidy:
	go mod tidy

swagger:
	swag init -g cmd/app/main.go --output docs

# ── docker ───────────────────────────────────────────────
docker-up:
	docker run --name todo-postgres \
		-e POSTGRES_USER=${POSTGRES_USER} \
		-e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
		-e POSTGRES_DB=${POSTGRES_DB} \
		-p ${POSTGRES_PORT}:5432 \
		-d postgres:16-alpine
	docker run --name todo-redis \
		-p 6379:6379 \
		-d redis:7-alpine

docker-down:
	docker stop todo-postgres todo-redis || true
	docker rm todo-postgres todo-redis || true

ps:
	docker ps

compose-up:
	docker compose up --build -d

compose-down:
	docker compose down

compose-logs:
	docker compose logs -f app

compose-fresh:
	docker compose down -v
	docker compose up --build -d
# ── migrations ───────────────────────────────────────────
migrate:
	docker exec -i todo-postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} \
		< migrations/001_create_todos.sql
	docker exec -i todo-postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} \
		< migrations/002_create_users.sql

migrate-up:
	migrate -path migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" down 1

migrate-version:
	migrate -path migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" version
# ── api ──────────────────────────────────────────────────
register:
	curl -s -X POST localhost:8080/v1/auth/register \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"securepassword"}' | jq

login:
	curl -s -X POST localhost:8080/v1/auth/login \
		-H "Content-Type: application/json" \
		-d '{"email":"test@example.com","password":"securepassword"}' | jq

create:
	curl -s -X POST localhost:8080/v1/todos \
		-H "Content-Type: application/json" \
		-H "Authorization: Bearer ${TOKEN}" \
		-d '{"title":"First authenticated todo"}' | jq

list:
	curl -s localhost:8080/v1/todos \
		-H "Authorization: Bearer ${TOKEN}" | jq

fail:
	curl -s localhost:8080/v1/todos | jq

verify:
	docker exec -it todo-redis redis-cli keys "*"
