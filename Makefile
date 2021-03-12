PROJECT_NAME := "fastdns"

.PHONY: all lint dep test bench race msan dep build mocks clean help

all: build

lint: ## Lint the files
	@golint -set_exit_status ./...

test: ## Run unit tests
	@go test -short -cover -v -count=1 ./...

bench: ## Run benchmark
	@go test -bench ./...

race: dep ## Run data race detector
	@go test -race ./...

msan: dep ## Run memory sanitizer
	@go test -msan ./...

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -v ./...

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
