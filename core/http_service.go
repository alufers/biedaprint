package core

import (
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

/*
HTTPService is responsible for serving the dashboard and the graphql endpoint.
*/
type HTTPService struct {
	app    *App
	router *gin.Engine
}

/*
NewHTTPService constructs a new HTTPService.
*/
func NewHTTPService(app *App) *HTTPService {
	return &HTTPService{
		app: app,
	}
}

/*
RunHTTPServer starts the http server using the router. It should be called in a spearate goroutine,
*/
func (hs *HTTPService) RunHTTPServer() {
	// Setting up Gin
	hs.router = gin.Default()
	hs.router.RedirectTrailingSlash = false
	hs.router.Use(cors.Default())
	hs.router.Any("/query", hs.graphqlQueryHandler())
	hs.router.GET("/playground", hs.graphqlPlaygroundHandler())
	hs.router.NoRoute(hs.frontendHandler())
	hs.router.Run(":4444")
}

func (hs *HTTPService) graphqlQueryHandler() gin.HandlerFunc {
	h := handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{
		App: hs.app,
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
func (hs *HTTPService) graphqlPlaygroundHandler() gin.HandlerFunc {
	h := handler.Playground("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
