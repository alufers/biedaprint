
MAKEFLAGS += -j2

build backend: build-frontend backend-graphql-codegen
	packr build -tags frontend_packr -o build/biedaprint server/server.go

backend-graphql-codegen:
	go run github.com/99designs/gqlgen
build-frontend: frontend-graphql-codegen
	cd frontend; \
	VUE_APP_GRAPHQL_HTTP=/query VUE_APP_GRAPHQL_WS_AUTO_RELATIVE=true node_modules/.bin/vue-cli-service build --dest ../static 

frontend-graphql-codegen:
	cd frontend; \
	npm run graphql-codegen