GO        ?= go
MOCKGEN   ?=mockgen

# set default shell
SHELL = bash -e -o pipefail

# Variables
VERSION                  ?= $(shell cat ./VERSION)

## Docker related
DOCKER_EXTRA_ARGS        ?=
DOCKER_REGISTRY          ?=
DOCKER_REPOSITORY        ?=
DOCKER_TAG               ?= ${VERSION}
IMAGE_NAME               := ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}case-study:${DOCKER_TAG}
DOCKER_BUILD_ARGS        ?=${DOCKER_EXTRA_ARGS} --build-arg version="${VERSION}"
GOLANGCI_LINT_BIN_PATH   ?= bin/golangci-lint

.PHONY: build test

default: build

help:
	@echo "Usage: make [<target>]"
	@echo "where available targets are:"
	@echo
	@echo "build                : Build PADM and PADM-DIAG binary"
	@echo "build-docker         : Build PADM and PADM-DIAG docker image"
	@echo "help                 : Print this help"
	@echo "test                 : Run unit tests, if any"
	@echo "dmi-mock             : Generates mock for DMI"
	@echo

build:
	mkdir -p bin
	$(GO) build -race -o bin/case-study \
	    -ldflags \
            "-X github.com/internal/pkg/cli.Version=${VERSION}" \
	    cmd/main.go

build-docker:
	docker build $(DOCKER_BUILD_ARGS) -t ${IMAGE_NAME} -f build/docker/case-study.dockerfile .

test:
	$(GO) test -race -v ./... -coverprofile coverage.out 2>&1 | tee output.log
