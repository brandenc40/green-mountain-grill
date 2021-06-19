BINARY_NAME=gmg

.PHONY: clean
clean:
	@rm ${BINARY_NAME}

.PHONY: build
build:
	@go build -o ${BINARY_NAME} main.go

.PHONY: test
test:
	@go test -v main.go

.PHONY: run
run: build
	@./${BINARY_NAME}

.PHONY: run-dev
run-dev: build
	@ENVIRONMENT=development ./${BINARY_NAME}
