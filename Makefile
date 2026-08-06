COMPOSE = docker compose -f deploy/docker-compose.yml --env-file .env.example
PYTHON ?= python

.PHONY: help up down stop reset logs ps rebuild-backend api-test api-check web-test web-check fmt fmt-check tidy ci-local lambda-zip compose-config check-docs check-static

help:
	@echo "PhishGuard simulation platform — common targets:"
	@echo "  up down stop reset logs ps rebuild-backend"
	@echo "  api-test api-check web-test web-check fmt fmt-check tidy"
	@echo "  ci-local lambda-zip compose-config check-docs check-static"

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

stop:
	$(COMPOSE) stop

reset:
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

api-check:
	cd backend && test -z "$$(gofmt -l .)" && go vet ./... && go test ./...

web-test:
	cd apps/web && npm test

web-check:
	cd apps/web && npm run check && npm test

fmt:
	cd backend && gofmt -w .
	cd apps/web && npm run format

fmt-check:
	cd backend && test -z "$$(gofmt -l .)"
	cd apps/web && npm run format:check

tidy:
	cd backend && go mod tidy

lambda-zip:
	cd backend && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/bootstrap ./cmd/lambda

compose-config:
	$(COMPOSE) config -q

# Stdlib Python 3.8+ scripts; PYTHON defaults to `python` (works on Windows too).
check-docs:
	$(PYTHON) scripts/check_links.py
	$(PYTHON) scripts/check_env_example.py

check-static: check-docs
	$(PYTHON) scripts/check_metrics_contract.py

ci-local: api-check web-check
