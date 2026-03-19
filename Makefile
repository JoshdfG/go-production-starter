-include .env
run: 
	go run ./cmd/app

d-run:
	docker run --name todo-postgres \
  -e POSTGRES_USER=${POSTGRES_USER} \
  -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} \
  -e POSTGRES_DB=${POSTGRES_DB} \
  -p ${POSTGRES_PORT}:5432 \
  -d postgres:16-alpine

redis:
	docker run --name todo-redis \
  -p 6379:6379 \
  -d redis:7-alpine

apply:
	docker exec -i todo-postgres psql -U ${POSTGRES_USER} -d ${POSTGRES_DB} < migrations/002_create_users.sql

ps:
	docker ps

test:
	go test ./internal/usecase/...

test-v:
	go test ./internal/usecase/... -v

register:
	curl -X POST localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"securepassword"}'

login:
	curl -X POST localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"securepassword"}'

# 3. Create a todo with auth
create:
	curl -X POST localhost:8080/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d '{"title":"First authenticated todo"}'

# 4. List todos with auth
list:
	curl localhost:8080/v1/todos \
  -H "Authorization: Bearer ${TOKEN}"

# 5.  without token — should get 401
fail:
	curl localhost:8080/v1/todos

# First call — should hit Postgres
1:
	curl localhost:8080/v1/todos \
  -H "Authorization: Bearer ${TOKEN}"

# Second call — should hit Redis only
2:
	curl localhost:8080/v1/todos \
  -H "Authorization: Bearer ${TOKEN}"

verify:
	docker exec -it todo-redis redis-cli keys "*"
