package biedaprint

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (app *App) RunHTTPServer() {
	// Setting up Gin
	r := gin.Default()
	r.Use(cors.Default())
	r.Any("/query", app.graphqlQueryHandler())
	r.GET("/playground", app.graphqlPlaygroundHandler())
	r.NoRoute(app.frontendHandler())
	r.Run(":4444")
}

func (app *App) graphqlQueryHandler() gin.HandlerFunc {
	h := handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{
		App: app,
	}}), handler.WebsocketUpgrader(websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func (app *App) graphqlPlaygroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
