//+build !frontend_packr

package biedaprint

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/phayes/freeport"
)

// proxies all frontend requests to a running vue-cli instance

func (app *App) initFrontend() int {
	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("sh", "-c", fmt.Sprintf("npm run serve -- --port %v", port))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir, err = filepath.Abs("./frontend")
	if err != nil {
		panic(err)
	}
	go func() {
		err = cmd.Run()
		if err != nil {
			log.Printf("Failed to run frontend command %v", err)
		}
	}()
	return port
}

func (app *App) frontendHandler() gin.HandlerFunc {
	// //port := app.initFrontend()
	// url, err := url.Parse(fmt.Sprintf("http://localhost:%d", 8080))
	// if err != nil {
	// 	panic(err)
	// }
	// proxy := httputil.NewSingleHostReverseProxy(url)

	return func(c *gin.Context) {
		//proxy.ServeHTTP(c.Writer, c.Request)
		c.JSON(200, ":((")
	}
}
