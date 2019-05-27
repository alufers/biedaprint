//+build !frontend_packr

package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) frontendHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<style> * { font-family: Arial, sans-serif; } </style>
		<h1>This is just a backend server for biedaprint.</h1> 
		Please use the Vue development server to access the frontend. 
		When building for production the frontend will be embedded into the app.<br>
		<a href="/playground">Graphql playground</a>`))
	}
}
