APP_NAME := taskr
BUILD_DIR := bin
GO := go

.PHONY: help build run test fmt vet check clean

help:
	@printf "Available targets:\n"
	@printf "  build  Build the application\n"
	@printf "  run    Run the application\n"
	@printf "  check  Check application for potential errors"
	@printf "  test   Run tests\n"
	@printf "  fmt    Format Go source files\n"
	@printf "  vet    Run go vet\n"
	@printf "  clean  Remove build artifacts\n"

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(APP_NAME) .

run:
	$(GO) run .

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

check: fmt vet test

clean:
	rm -rf $(BUILD_DIR)
