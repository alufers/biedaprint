
MAKEFLAGS += -j2

build_with_out = packr build -tags frontend_packr -ldflags="-s -w" -o build/$(1) server/server.go

build-backend: build-frontend backend-graphql-codegen
	packr build -tags frontend_packr -ldflags="-s -w" -o build/biedaprint server/server.go

build-multiplatform: build-frontend backend-graphql-codegen
	GOOS=linux GOARCH=arm GOARM=7 $(call build_with_out,biedaprint-linux-armv7)
	GOOS=linux GOARCH=arm64 $(call build_with_out,biedaprint-linux-arm64)
	GOOS=linux GOARCH=amd64 $(call build_with_out,biedaprint-linux-amd64)
	GOOS=darwin GOARCH=amd64 $(call build_with_out,biedaprint-macos-amd64)

backend-graphql-codegen:
	go run github.com/99designs/gqlgen
build-frontend: frontend-graphql-codegen
	cd frontend; \
	VUE_APP_GRAPHQL_HTTP=/query VUE_APP_GRAPHQL_WS_AUTO_RELATIVE=true node_modules/.bin/vue-cli-service build --dest ../static 

frontend-graphql-codegen:
	cd frontend; \
	npm run graphql-codegen