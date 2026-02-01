package app

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func (a App) runSkills(ctx context.Context, args []string) int {
	if len(args) == 0 || isHelp(args[0]) {
		a.Out.Raw(`gkn skills <command>

Commands:
  clone [--remote url]
  sync [--target name] [--mode copy|mirror|link]
  watch [--target name] [--interval sec]
  diff [--target name]
  verify [--target name]
  status [--target name]
  link [--target name]
  pin --target name --ref <commit|tag>
  clean [--target name]

Common flags:
  --only <glob> (repeatable)
  --exclude <glob> (repeatable)
  --dry-run
  --force`)
		return 0
	}
	switch args[0] {
	case "clone":
		return a.runSkillsClone(ctx, args[1:])
	case "sync":
		return a.runSkillsSync(ctx, args[1:])
	case "watch":
		return a.runSkillsWatch(ctx, args[1:])
	case "diff":
		return a.runSkillsDiff(ctx, args[1:])
	case "verify":
		return a.runSkillsVerify(ctx, args[1:])
	case "status":
		return a.runSkillsStatus(ctx, args[1:])
	case "link":
		return a.runSkillsLink(ctx, args[1:])
	case "pin":
		return a.runSkillsPin(ctx, args[1:])
	case "clean":
		return a.runSkillsClean(ctx, args[1:])
	case "--help", "-h", "help":
		fs := flag.NewFlagSet("skills", flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		_ = fs.Parse([]string{})
		a.Out.Raw("use 'gkn skills <command> --help' for details")
		return 0
	default:
		a.Out.Err(fmt.Sprintf("unknown skills command: %s", args[0]), nil)
		return 1
	}
}
