_gkn_complete() {
  local cur prev
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"

  if [[ ${COMP_CWORD} -eq 1 ]]; then
    COMPREPLY=( $(compgen -W "help repo skills config doctor version clone quickstart" -- "$cur") )
    return 0
  fi

  case "$prev" in
    repo)
      COMPREPLY=( $(compgen -W "list status open path recent info graph clone exec" -- "$cur") )
      return 0
      ;;
    skills)
      COMPREPLY=( $(compgen -W "clone sync link watch diff verify status pin clean" -- "$cur") )
      return 0
      ;;
    config)
      COMPREPLY=( $(compgen -W "show init validate" -- "$cur") )
      return 0
      ;;
  esac
}

complete -F _gkn_complete gkn
