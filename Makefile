compose-up:
	docker-compose up -d

docs-generate:
	swag init -g ./cmd/main.go


build:
	go build ./cmd/main.go

run:
	go run ./cmd/main.go

test:
	go test ./tests -v

test-services:
	go test ./internal/services_test/ -v

clean-test-cache:
	go clean -testcache

lint:
	golangci-lint run -v

overlord-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				./pkg/overlord/overlord.proto