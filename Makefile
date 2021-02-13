VERSION := $(shell cat ./VERSION)
BIN_NAME := gomd5sum

.PHONY: build
build:
	go build \
		-o ./bin/$(BIN_NAME) \
		-ldflags "-X main.version=$(VERSION) -X main.name=$(BIN_NAME)" \
		.

.PHONY: integrationtest
integrationtest: build
	./integrationtest.sh
