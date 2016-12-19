.PHONY: all microkit-example



GIT_COMMIT := $(shell git rev-parse --short HEAD)



all: microkit-example

microkit-example:
	@go build \
		-o microkit-example \
		-ldflags " \
			-X main.gitCommit=$(GIT_COMMIT) \
		" \
		.
