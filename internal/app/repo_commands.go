package app

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func (a App) runRepo(ctx context.Context, args []string) int {
	if len(args) == 0 || args[0] == "help" {
		a.Out.Raw(`gkn repo <command>

Commands:
  list
  status
  cd <pattern> [--pick n]
  open <pattern> [--pick n]
  path <pattern> [--pick n]
  recent [--limit n]
  info <pattern> [--pick n]
  graph <pattern> [--pick n] [--limit n]
  clone <url> [--name repo]
  exec --cmd "<command>" [--parallel n] [--timeout sec] [--require-clean]

Common flags:
  --only <glob> (repeatable)
  --exclude <glob> (repeatable)
  --dry-run
  --force`)
		return 0
	}
	switch args[0] {
	case "list":
		return a.runRepoList(ctx, args[1:])
	case "status":
		return a.runRepoStatus(ctx, args[1:])
	case "cd":
		return a.runRepoCd(ctx, args[1:])
	case "open":
		return a.runRepoOpen(ctx, args[1:])
	case "path":
		return a.runRepoPath(ctx, args[1:])
	case "recent":
		return a.runRepoRecent(ctx, args[1:])
	case "info":
		return a.runRepoInfo(ctx, args[1:])
	case "graph":
		return a.runRepoGraph(ctx, args[1:])
	case "clone":
		return a.runRepoClone(ctx, args[1:])
	case "exec":
		return a.runRepoExec(ctx, args[1:])
	case "--help", "-h":
		fs := flag.NewFlagSet("repo", flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		_ = fs.Parse([]string{})
		a.Out.Raw("use 'gkn repo <command> --help' for details")
		return 0
	default:
		a.Out.Err(fmt.Sprintf("unknown repo command: %s", args[0]), nil)
		return 1
	}
}
