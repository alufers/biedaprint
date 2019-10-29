
MAKEFLAGS += -j2

PACKR2_PATH ?= packr2


ifeq ($(TRAVIS_TAG),)
TRAVIS_TAG := $(GITHUB_REF:refs/tags/%=%)
endif

ifeq ($(TRAVIS_TAG),)
TRAVIS_TAG := "development"
endif

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
build_with_out = CGO_ENABLED=1 go build -tags frontend_packr -ldflags="-s -w -X github.com/alufers/biedaprint/core.AppVersion=$(TRAVIS_TAG) -X github.com/alufers/biedaprint/core.AppReleaseExecutableName=$(1)" -o $(ROOT_DIR)/build/$(1) server/server.go

build-backend: build-frontend backend-graphql-codegen
	$(PACKR2_PATH) build -tags frontend_packr -ldflags="-s -w" -o $(ROOT_DIR)/build/biedaprint server/server.go

build-multiplatform: backend-graphql-codegen embed-assets
	GOOS=linux GOARCH=arm GOARM=7 CC_FOR_TARGET=arm-linux-gnueabi-gcc CC=arm-linux-gnueabi-gcc $(call build_with_out,biedaprint-linux-armv7)
	GOOS=linux GOARCH=arm64 CC_FOR_TARGET=aarch64-linux-gnu-gcc CC=aarch64-linux-gnu-gcc $(call build_with_out,biedaprint-linux-arm64)
	GOOS=linux GOARCH=amd64 $(call build_with_out,biedaprint-linux-amd64)
	# disable darwin for now, too much time to compile
	# GOOS=darwin GOARCH=amd64 $(call build_with_out,biedaprint-macos-amd64) 


embed-assets: build-frontend
	cd core; \
	$(PACKR2_PATH) -v

backend-graphql-codegen:
	# go run github.com/99designs/gqlgen

build-frontend: frontend-graphql-codegen
	cd frontend; \
	VUE_APP_GRAPHQL_HTTP=/query VUE_APP_GRAPHQL_WS_AUTO_RELATIVE=true node_modules/.bin/vue-cli-service build --dest ../static 

frontend-graphql-codegen:
	cd frontend; \
	npm run graphql-codegen
