BINARY_NAME=gmg

.PHONY: default
default: clean generate build

.PHONY: build
build:
	@go build -o ${BINARY_NAME} main.go

.PHONY: clean
clean:
	@rm ${BINARY_NAME}

.PHONY: test
test:
	@go test ./... -coverprofile cp.out

.PHONY: coverage
coverage: test
	@go tool cover -html=cp.out

.PHONY: run-prod
run-prod: build
	@echo "Running in Production mode"
	@ENVIRONMENT=production ./${BINARY_NAME}

.PHONY: run
run: build
	@./${BINARY_NAME}

.PHONY: generate
generate: .gen .tidy

.PHONY: .gen
.gen:
	@go get github.com/alvaroloes/enumer
	@go generate -x ./...

.PHONY: .tidy
.tidy:
	@go mod tidy
