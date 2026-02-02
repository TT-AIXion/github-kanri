package main

import (
	"context"
	"os"

	"github.com/TT-AIXion/github-kanri/internal/app"
	"github.com/TT-AIXion/github-kanri/internal/output"
)

var Version = "dev"
var exitFunc = os.Exit

func main() {
	exitFunc(runMain(os.Args[1:]))
}

func runMain(args []string) int {
	ctx := context.Background()
	jsonMode, args := consumeJSONFlag(args)
	application := app.App{Version: Version, Out: output.New(jsonMode)}
	return application.Run(ctx, args)
}

func consumeJSONFlag(args []string) (bool, []string) {
	var out []string
	jsonMode := false
	for _, arg := range args {
		if arg == "--json" {
			jsonMode = true
			continue
		}
		out = append(out, arg)
	}
	return jsonMode, out
}
