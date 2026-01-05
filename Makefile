COMPOSE = docker compose -f deploy/docker-compose.yml --env-file .env.example

.PHONY: up down logs ps rebuild-backend api-test web-test fmt tidy ci-local lambda-zip

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f --tail=100

ps:
	$(COMPOSE) ps

# Rebuild API + worker after Go changes without restarting the whole stack.
rebuild-backend:
	$(COMPOSE) up --build -d api worker

api-test:
	cd backend && go test ./...

web-test:
	cd apps/web && npm test

fmt:
	cd backend && gofmt -w .
	cd apps/web && npm run format

tidy:
	cd backend && go mod tidy

lambda-zip:
	cd backend && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0&& go build -o bin/bootstrap ./cmd/lambda

ci-local: api-test
	cd apps/web && npm run check && npm test
