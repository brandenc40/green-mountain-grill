SERVER_BINARY_NAME=server
FRONTEND_PATH = frontend/

export APP_ROT = $(CURDIR)

.PHONY: default
default: build-js generate build mocks

.PHONY: build
build: build-js
	@echo "Building Binary..."
	@go build -o bin/${SERVER_BINARY_NAME} cmd/server/main.go

.PHONY: mocks
mocks:
	@echo "Building Mocks..."
	@mockery --all --dir ./server
	@mockery --name=Client

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
	@echo "Generating Code..."
	@go get github.com/alvaroloes/enumer
	@go mod vendor
	@go generate -x ./...

.PHONY: .tidy
.tidy:
	@go mod tidy

.PHONY: build-js
build-js:
	@echo Building React App..
	@npm run-script build -prefix $(FRONTEND_PATH)

.PHONY: dev-react
dev-react:
	@npm run-script start -prefix $(FRONTEND_PATH)

.PHONY: dev-go
dev-go:
	@go run cmd/server/main.go
