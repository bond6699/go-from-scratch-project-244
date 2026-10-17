.PHONY: run build

run:
	go run cmd/gendiff/main.go $(ARGS)
build:
	go build -o bin/gendiff ./cmd/gendiff
lint:
	golangci-lint run
lint-fix:
	golangci-lint run --fix
test:
	go test -cover -v ./...