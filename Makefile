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
SERVICE_PORT?=8080
DOCKER_REGISTRY?= #if set it should finished by /
EXPORT_RESULT?=false # for CI please set EXPORT_RESULT to true

.DEFAULT_GOAL := all

.PHONY: all test build vendor

all: clean vendor lint test build


## Generate:
codegen: ## Generate OpenAPI 3.0 client and server code
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

lint-go: 	## Use golintci-lint on your project
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--out-format checkstyle ./... | tee /dev/tty > checkstyle-report.xml" || echo "" ))
	docker run --name make-lint-go --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --verbose --deadline=120s $(OUTPUT_OPTIONS)

## Test:
test:		## Run the tests of the project
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off go get -u github.com/jstemmer/go-junit-report
	$(eval OUTPUT_OPTIONS = | tee /dev/tty | go-junit-report -set-exit-code > junit-report.xml)
endif
	$(GO) test -v -race ./... $(OUTPUT_OPTIONS)

coverage: 	## Run the tests of the project and export the coverage
	$(GO) test -cover -covermode=count -coverprofile=profile.cov ./...
	$(GO) tool cover -func profile.cov
ifeq ($(EXPORT_RESULT), true)
	GO111MODULE=off $(GO) get -u github.com/AlekSi/gocov-xml
	GO111MODULE=off $(GO) get -u github.com/axw/gocov/gocov
	gocov convert profile.cov | gocov-xml > coverage.xml
endif


## Build:
build: 		## Build your project and put the output binary in out/bin/
	mkdir -p ./bin
	GO111MODULE=on $(GO) build -mod vendor -o ./bin/$(BINARY_NAME) ./cmd/up2

clean: 		## Remove build related file
	rm -fr ./bin
	rm -f ./junit-report.xml checkstyle-report.xml ./coverage.xml ./profile.cov yamllint-checkstyle.xml

vendor: 	## Copy of all packages needed to support builds and tests in the vendor directory
	$(GO) mod vendor

watch: 		## Run the code with cosmtrek/air to have automatic reload on changes
	$(eval PACKAGE_NAME=$(shell head -n 1 go.mod | cut -d ' ' -f2))
	docker run --name make-watch -it --rm -w /go/src/$(PACKAGE_NAME) -v $(shell pwd):/go/src/$(PACKAGE_NAME) -p $(SERVICE_PORT):$(SERVICE_PORT) cosmtrek/air


## Docker:
docker-lint: ## Lint your Dockerfile (if present)
ifeq ($(shell test -e ./Dockerfile && echo yes),yes)
	$(eval CONFIG_OPTION = $(shell [ -e $(shell pwd)/.hadolint.yaml ] && echo "-v $(shell pwd)/.hadolint.yaml:/root/.config/hadolint.yaml" || echo "" ))
	$(eval OUTPUT_OPTIONS = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "--format checkstyle" || echo "" ))
	$(eval OUTPUT_FILE = $(shell [ "${EXPORT_RESULT}" == "true" ] && echo "| tee /dev/tty > checkstyle-report.xml" || echo "" ))
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