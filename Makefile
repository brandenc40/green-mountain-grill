SERVER_BINARY_NAME=server

.PHONY: default
default: generate build

.PHONY: build
build:
	@go build -o bin/${SERVER_BINARY_NAME} cmd/server/main.go

.PHONY: clean
clean:
	@rm bin/${SERVER_BINARY_NAME}

.PHONY: test
test:
	@go test ./... -coverprofile cp.out

.PHONY: coverage
coverage: test
	@go tool cover -html=cp.out

.PHONY: run-prod
run-prod: build
	@echo "Running in Production mode"
	@ENVIRONMENT=production bin/${SERVER_BINARY_NAME}

.PHONY: run
run: build
	@bin/${SERVER_BINARY_NAME}

.PHONY: generate
generate: .gen .tidy

.PHONY: .gen
.gen:
	@go get github.com/alvaroloes/enumer
	@go mod vendor
	@go generate -x ./...

.PHONY: .tidy
.tidy:
	@go mod tidy
