// +build frontend_packr

package core

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

// http interceptor for interceping 404 and serving index.html so that the SPA works

type interceptResponseWriter struct {
	http.ResponseWriter
	errH func(http.ResponseWriter, int)
}

func (w *interceptResponseWriter) WriteHeader(status int) {
	if status >= http.StatusBadRequest || status == 301 {
		w.errH(w.ResponseWriter, status)
		w.errH = nil
	} else {
		w.ResponseWriter.WriteHeader(status)
	}
}

//ErrorHandler handles intercepted errors
type ErrorHandler func(http.ResponseWriter, int)

func (w *interceptResponseWriter) Write(p []byte) (n int, err error) {
	if w.errH == nil {
		return len(p), nil
	}
	return w.ResponseWriter.Write(p)
}

func interceptHandler(next http.Handler, errH ErrorHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(&interceptResponseWriter{w, errH}, r)
	})
}

func (hs *HTTPService) frontendHandler() gin.HandlerFunc {
	box := packr.New("Static frontend box", "../static")
	fs := interceptHandler(http.FileServer(box), func(w http.ResponseWriter, status int) {
		data, _ := box.FindString("index.html")
		w.Header().Set("Content-type", "text/html")
		fmt.Fprint(w, data)
	})
	return func(c *gin.Context) {
		fs.ServeHTTP(c.Writer, c.Request)

	}
}
