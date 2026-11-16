COVERAGE_MIN ?= 80

.PHONY: run build lint lint-fix test test-coverage test-with-cover pre-commit

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
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out | awk -v min=$(COVERAGE_MIN) '/total:/ { \
		gsub("%", "", $$NF); \
		if ($$NF < min) { \
			printf "FAIL: Coverage %.1f%% is below threshold %s%%\n", $$NF, min; \
			exit 1; \
		} else { \
			printf "PASS: Coverage %.1f%% meets threshold %s%%\n", $$NF, min; \
		} \
	}'
test-with-cover:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

pre-commit:
	make lint-fix
	make test-coverage

	