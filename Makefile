# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help


title:
	@echo "lucha Makefile"
	@echo "--------------"

build: ## Builds the application
	go build

test: ## Runs tests and coverage
	go test -v -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

install: build ## Builds an executable local version of lucha and puts in in /usr/local/bin
	sudo chmod +x lucha
	sudo mv lucha /usr/local/bin

run: build #Runs the default recursive lucha binary
	./lucha scan --rules-file lucha.yaml --recursive .

all: title build test ## Makes all targets

