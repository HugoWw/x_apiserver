package main

import (
	"fmt"
	"github.com/HugoWw/x_apiserver/cmd/x_apiserver/app"
	"os"
)

func main() {
	cmd := app.NewAPIServerCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, "X_ApiServer start up fail: ", err.Error())
		os.Exit(1)
	}
}
