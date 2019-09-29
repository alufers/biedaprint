
MAKEFLAGS += -j2

ifeq ($(TRAVIS_TAG),)
TRAVIS_TAG := "development"
endif

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
build_with_out = GO111MODULE=on go run github.com/gobuffalo/packr/v2/packr2 build -tags frontend_packr -ldflags="-s -w -X github.com/alufers/biedaprint/core.AppVersion=$(TRAVIS_TAG) -X github.com/alufers/biedaprint/core.AppReleaseExecutableName=$(1)" -o $(ROOT_DIR)/build/$(1) server/server.go

build-backend: build-frontend backend-graphql-codegen
	go run github.com/gobuffalo/packr/v2/packr2 build -tags frontend_packr -ldflags="-s -w" -o $(ROOT_DIR)/build/biedaprint server/server.go

build-multiplatform: build-frontend backend-graphql-codegen
	GOOS=linux GOARCH=arm GOARM=7 $(call build_with_out,biedaprint-linux-armv7)
	GOOS=linux GOARCH=arm64 $(call build_with_out,biedaprint-linux-arm64)
	GOOS=linux GOARCH=amd64 $(call build_with_out,biedaprint-linux-amd64)
	GOOS=darwin GOARCH=amd64 $(call build_with_out,biedaprint-macos-amd64)

backend-graphql-codegen:
	# go run github.com/99designs/gqlgen

build-frontend: frontend-graphql-codegen
	cd frontend; \
	VUE_APP_GRAPHQL_HTTP=/query VUE_APP_GRAPHQL_WS_AUTO_RELATIVE=true node_modules/.bin/vue-cli-service build --dest ../static 

frontend-graphql-codegen:
	cd frontend; \
	npm run graphql-codegen
