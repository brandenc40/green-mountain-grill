BINARY_NAME=gmg

.PHONY: clean
clean:
	@rm ${BINARY_NAME}

.PHONY: build
build:
	@go build -o ${BINARY_NAME} main.go

.PHONY: test
test:
	@go test ./... -coverprofile cp.out

.PHONY: coverage
coverage: test
	@go tool cover -html=cp.out

.PHONY: run
run: build
	@./${BINARY_NAME}

.PHONY: run-dev
run-dev: build
	@ENVIRONMENT=development ./${BINARY_NAME}
