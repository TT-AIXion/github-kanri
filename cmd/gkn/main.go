package main

import (
	"context"
	"os"

	"github.com/AIXion-Team/github-kanri/internal/app"
	"github.com/AIXion-Team/github-kanri/internal/output"
)

var Version = "dev"

func main() {
	ctx := context.Background()
	jsonMode, args := consumeJSONFlag(os.Args[1:])
	application := app.App{Version: Version, Out: output.New(jsonMode)}
	os.Exit(application.Run(ctx, args))
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
