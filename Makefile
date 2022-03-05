ifneq (, $(shell which tput))
	GREEN  := $(shell tput -Txterm setaf 2)
	YELLOW := $(shell tput -Txterm setaf 3)
	WHITE  := $(shell tput -Txterm setaf 7)
	CYAN   := $(shell tput -Txterm setaf 6)
	RESET  := $(shell tput -Txterm sgr0)
endif

GO=go
BINARY_NAME=up2
VERSION?=1.0.0
DOCKER_REGISTRY?= #if set it should finished by /
CI?=false

.DEFAULT_GOAL := all

.PHONY: all test build vendor

all: clean vendor lint build test


## Generate:
codegen: ## Generate OpenAPI 3.0 server code
ifeq (, $(shell which oapi-codegen))
	GO111MODULE=off $(GO) get -u github.com/deepmap/oapi-codegen/cmd/oapi-codegen
endif
	oapi-codegen -package server -generate types api/spec.yaml > pkg/server/types.gen.go
	oapi-codegen -package server -generate server api/spec.yaml > pkg/server/server.gen.go


## Format:
format:	## Format source files (alias: fmt)
	$(GO) fmt ./...

fmt: format


## Lint:
lint:		## Run all available linters
lint: lint-go

lint-go: 	## Use golangci-lint on your project
	$(eval OUTPUT_OPTIONS = $(shell [ "${CI}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	golangci-lint run --verbose $(OUTPUT_OPTIONS)

## Test:
test:		## Run the unit and integration tests
ifeq ($(CI), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GO) test --tags=integration -v -race  ./... $(OUTPUT_OPTIONS)

test-cov: 		## Run the unit and integration tests and export the coverage
	$(GO) test --tags=integration -cover -covermode=atomic -coverprofile=out/coverage.out ./...
	$(GO) tool cover -func out/coverage.out

unittest:		## Run the unit tests
ifeq ($(CI), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GO) test -v -race ./... $(OUTPUT_OPTIONS)

unittest-cov: 	## Run the unit tests and export the coverage
	$(GO) test -cover -covermode=atomic -coverprofile=out/coverage.out ./...
	$(GO) tool cover -func out/coverage.out


## Build:
build: 		## Build your project and put the output binary in out/bin/
	mkdir -p ./out
	GO111MODULE=on $(GO) build -mod vendor -o ./out/bin/$(BINARY_NAME) ./cmd/up2

clean: 		## Remove build related file
	rm -fr ./out ./tmp
	rm -f ./junit-report.xml checkstyle-report.xml yamllint-checkstyle.xml

vendor: 	## Copy of all packages needed to support builds and tests in the vendor directory
	$(GO) mod vendor


## Docker:
docker-lint: ## Lint your Dockerfile (if present)
ifeq ($(shell test -e ./Dockerfile && echo yes),yes)
	$(eval CONFIG_OPTION = $(shell [ -e $(shell pwd)/.hadolint.yaml ] && echo "-v $(shell pwd)/.hadolint.yaml:/root/.config/hadolint.yaml" || echo "" ))
	$(eval OUTPUT_OPTIONS = $(shell [ "${CI}" == "true" ] && echo "--format checkstyle" || echo "" ))
	$(eval OUTPUT_FILE = $(shell [ "${CI}" == "true" ] && echo "| tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --name make-docker-lint --rm -i $(CONFIG_OPTION)  hadolint/hadolint hadolint $(OUTPUT_OPTIONS) - < ./Dockerfile $(OUTPUT_FILE)
endif

docker-build: ## Use the dockerfile to build the container
	docker build --rm --tag $(BINARY_NAME) .

docker-release: ## Release the container with tag latest and version
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker tag $(BINARY_NAME) $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):latest
	docker push $(DOCKER_REGISTRY)$(BINARY_NAME):$(VERSION)


## Help:
help:		## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)