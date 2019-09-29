package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/alufers/biedaprint/core"
)

func main() {
	var runInBackground bool

	flag.BoolVar(&runInBackground, "b", false, "run the app in background (disown)")
	flag.Parse()
	if runInBackground {
		bin, err := os.Executable()
		if err != nil {
			panic(fmt.Sprintf("failed to get current executable: %v", err))
		}
		cmd := exec.Command(bin)
		err = cmd.Start()
		if err != nil {
			panic(err)
		}
		return
	}
	app := core.NewApp()
	app.Init()
	app.RunHTTPServer()
}
