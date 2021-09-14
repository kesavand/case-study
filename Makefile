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
EVNT_PRODCUER_IMAGE_NAME := ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}evnt-producer:${DOCKER_TAG}
EVNT_RCVR_IMAGE_NAME     := ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}evnt-rcvr:${DOCKER_TAG}
DOCKER_BUILD_ARGS        ?=${DOCKER_EXTRA_ARGS} --build-arg version="${VERSION}"
GOLANGCI_LINT_BIN_PATH   ?= bin/golangci-lint

.PHONY: build test

default: build


build:
	mkdir -p bin
	$(GO) build -race -o bin/evnt-producer \
	    -ldflags \
            "-X github.com/internal/pkg/cli.Version=${VERSION}" \
	    cmd/evnt-rcvr/main.go
	$(GO) build -race -o bin/evnt-rcvr \
	    -ldflags \
            "-X github.com/internal/pkg/cli.Version=${VERSION}" \
	    cmd/evnt-producer/main.go

build-docker:
	docker build $(DOCKER_BUILD_ARGS) -t ${EVNT_PRODCUER_IMAGE_NAME} -f build/docker/evnt-producer.dockerfile .
	docker build $(DOCKER_BUILD_ARGS) -t ${EVNT_RCVR_IMAGE_NAME} -f build/docker/evnt-rcvr.dockerfile .

test:
	$(GO) test -race -v ./... -coverprofile coverage.out 2>&1 | tee output.log
