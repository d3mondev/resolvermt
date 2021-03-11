PROJECT_NAME := "fastdns"
PKG := "github.com/d3mondev/fastdns"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)

.PHONY: all lint dep test race msan dep build mocks clean help

all: build

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unit tests
	@go test -cover -v ${PKG_LIST}

race: dep ## Run data race detector
	@go test -race ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@go test -msan ${PKG_LIST}

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -v $(PKG)

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
