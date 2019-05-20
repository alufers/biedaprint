package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/alufers/biedaprint"
)

const defaultPort = "4444"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	app := biedaprint.NewApp()
	http.Handle("/query", handler.GraphQL(biedaprint.NewExecutableSchema(biedaprint.Config{Resolvers: &biedaprint.Resolver{
		App: app,
	}})))

	app.Init()

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
