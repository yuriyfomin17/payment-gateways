.PHONY: dc run test lint

dc:
	docker-compose up -d  --remove-orphans --build

build:
	go build -race -o app cmd/main.go

run:
	go build -race -o app cmd/main.go && \
	HTTP_ADDR=:8080 \
	DSN="postgres://user:password@localhost:5432/payments?sslmode=disable" \
	MIGRATIONS_PATH="file://./internal/app/migrations" \
	RABBITMQ_URL="amqp://guest:guest@localhost:5672/" \
	REDIS_URL="localhost:6379" \
	./app

test:
	go test -race ./...

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2

lint:
	golangci-lint run ./...

generate:
	go generate ./...