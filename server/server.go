package main

import (
	"github.com/alufers/biedaprint/core"
)

func main() {
	app := core.NewApp()
	app.Init()
	app.RunHTTPServer()
}
