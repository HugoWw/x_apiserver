package main

import (
	"fmt"
	"github.com/HugoWw/x_apiserver/cmd/x_apiserver/app"
	"testing"
)

func TestAPIServer(t *testing.T) {
	cmd := app.NewAPIServerCommand()
	args := []string{
		"--bind-addr=0.0.0.0:8866",
		"--ctrl-conf=https://192.168.100.1:31803",
	}

	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		fmt.Printf("run x_apiserver error:%v\n", err)
	}
}
