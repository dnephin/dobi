package main

import (
	"github.com/dnephin/dobi/cmd"
	"github.com/dnephin/dobi/logging"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		logging.Log.Fatal(err)
	}
}
