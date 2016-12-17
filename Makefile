.PHONY: all microkit-example



GIT_COMMIT := $(shell git rev-parse --short HEAD)
PROJECT_VERSION=$(shell cat VERSION)



all: microkit-example

microkit-example:
	@go build \
		-o microkit-example \
		-ldflags " \
			-X main.gitCommit=$(GIT_COMMIT) \
			-X main.projectVersion=$(PROJECT_VERSION) \
		" \
		.
