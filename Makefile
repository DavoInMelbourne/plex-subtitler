SHELL:=/bin/bash

# project details
PRODUCT = papi
APPNAME = api
PACKAGE = github.com/cultureamp/public-api/api

# build variables
BRANCH_NAME ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE  ?= $(shell date '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT  ?= $(shell git rev-list -1 HEAD)
VERSION     ?= 0.0.0

BUILD_OVERRIDES = \
	-X "$(PACKAGE)/pkg/app.Name=$(APPNAME)" \
	-X "$(PACKAGE)/pkg/app.Product=$(PRODUCT)" \
	-X "$(PACKAGE)/pkg/app.Branch=$(BRANCH_NAME)" \
	-X "$(PACKAGE)/pkg/app.BuildDate=$(BUILD_DATE)" \
	-X "$(PACKAGE)/pkg/app.Commit=$(GIT_COMMIT)" \
	-X "$(PACKAGE)/pkg/app.Version=$(VERSION)" \

CMDPATH = ./cmd
DIST_DIR = ./dist/

.PHONY: format
format:
	go fmt ./internal/... ./pkg/... ./test/...

.PHONY: clean-dist
clean-dist:
	rm -rf $(DIST_DIR)
	mkdir -p $(DIST_DIR)

.PHONY: build-main
build-main:
	GOOS=linux GOARCH=amd64 go build -a \
		-ldflags='-w -s $(BUILD_OVERRIDES)' \
		-o $(DIST_DIR) $(CMDPATH)/main.go

.PHONY: build-api
build-api: clean-dist build-main
