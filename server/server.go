package main

import (
	"github.com/alufers/biedaprint"
)

func main() {
	app := biedaprint.NewApp()
	app.Init()
	app.RunHTTPServer()
}
