# Copyright (c) 2021 Kross IAM Project.
# https://github.com/krossdev/iam-ms/blob/main/LICENSE
#
LDFLAGS ?= -s -w

TAGS ?=

all: build

# build a production executable
build:
	go build -ldflags='$(LDFLAGS)' -tags '${TAGS}' -o kiam-ms

# startup development server with default args
dev: build
	./kiam-ms -c .config/config.yaml -watch

# startup development server with customized args
run:
	go run -ldflags='$(LDFLAGS)' -tags '${TAGS}' . -c .config/config.yaml -watch

# run test
test:
	go test $(arg) ./test/...

# clean
clean:
	go clean

.PHONY: build dev test all clean
