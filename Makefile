APP_NAME=relay-engine

run:
	go run cmd/api/main.go

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker ps

include .env
export

migrate-up:
	migrate \
	-path migrations \
	-database "$(DATABASE_URL)" \
	up

migrate-down:
	migrate \
	-path migrations \
	-database "$(DATABASE_URL)" \
	down 1

migrate-force:
	migrate \
	-path migrations \
	-database "$(DATABASE_URL)" \
	force 1

create-migration:
	migrate create -ext sql -dir migrations -seq $(name)