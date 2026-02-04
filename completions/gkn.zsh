#compdef gkn

local -a commands
commands=(
  'help:show help'
  'cd:print repo path for cd'
  'shell:shell integration'
  'repo:repo operations'
  'skills:skills sync operations'
  'config:config operations'
  'doctor:environment checks'
  'version:show version'
  'clone:clone repo'
  'quickstart:create repo'
)

if (( CURRENT == 2 )); then
  _describe -t commands command commands
  return
fi

case "$words[2]" in
  repo)
    local -a repo_cmds
    repo_cmds=(
      'list:list repos'
      'status:repo status'
      'cd:print repo path for cd'
      'open:open repo'
      'path:repo path'
      'recent:recent repos'
      'info:repo info'
      'graph:repo graph'
      'clone:clone repo'
      'exec:exec command'
    )
    _describe -t commands command repo_cmds
    ;;
  shell)
    local -a shell_cmds
    shell_cmds=(
      'install:install shell function'
      'zsh:print zsh function'
      'bash:print bash function'
      'fish:print fish function'
      'powershell:print powershell function'
    )
    _describe -t commands command shell_cmds
    ;;
  skills)
    local -a skills_cmds
    skills_cmds=(
      'clone:clone skills'
      'sync:sync skills'
      'link:link skills'
      'watch:watch skills'
      'diff:diff skills'
      'verify:verify skills'
      'status:skills status'
      'pin:pin skills'
      'clean:clean skills'
    )
    _describe -t commands command skills_cmds
    ;;
  config)
    local -a config_cmds
    config_cmds=(
      'show:show config'
      'init:init config'
      'validate:validate config'
    )
    _describe -t commands command config_cmds
    ;;
  *)
    ;;
esac
