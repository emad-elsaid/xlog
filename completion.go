package xlog

import (
	"fmt"
	"os"
)

// handleCompletion outputs completion script for requested shell.
func handleCompletion(shell string) {
	switch shell {
	case "bash":
		fmt.Print(bashCompletionScript)
	case "zsh":
		fmt.Print(zshCompletionScript)
	case "fish":
		fmt.Print(fishCompletionScript)
	default:
		fmt.Fprintf(os.Stderr, "Unknown shell: %s. Supported: bash, zsh, fish\n", shell)
		osExit(1)
	}
}

const bashCompletionScript = `# bash completion for xlog

_xlog_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # All flags from xlog -h
    opts="-source -build -sitename -index -notfoundpage -readonly -bind -serve-insecure -csrf-cookie -disabled-extensions -codestyle -theme -completion"

    # Flag-specific completions
    case "${prev}" in
        -theme)
            COMPREPLY=( $(compgen -W "light dark" -- ${cur}) )
            return 0
            ;;
        -completion)
            COMPREPLY=( $(compgen -W "bash zsh fish" -- ${cur}) )
            return 0
            ;;
        -source|-build)
            COMPREPLY=( $(compgen -d -- ${cur}) )
            return 0
            ;;
        -disabled-extensions)
            COMPREPLY=( $(compgen -W "all" -- ${cur}) )
            return 0
            ;;
    esac

    # Complete flags
    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi
}

complete -F _xlog_completion xlog

# Installation:
# Save to ~/.bash_completion.d/xlog or add to ~/.bashrc:
# eval "$(xlog -completion bash)"
`

const zshCompletionScript = `#compdef xlog

# zsh completion for xlog

_xlog() {
    local -a flags
    flags=(
        '-source[Directory that will act as a storage]:directory:_directories'
        '-build[Build all pages as static site in this directory]:directory:_directories'
        '-sitename[Site name]:name:'
        '-index[Index file name used as home page]:filename:'
        '-notfoundpage[Custom not found page]:filename:'
        '-readonly[Should xlog hide write operations]'
        '-bind[IP and port to bind the web server to]:address:'
        '-serve-insecure[Accept http connections]'
        '-csrf-cookie[CSRF cookie name]:name:'
        '-disabled-extensions[disable list of extensions]:extensions:(all)'
        '-codestyle[code highlighting style name]:style:'
        '-theme[bulma theme to use]:theme:(light dark)'
        '-completion[Generate shell completion script]:shell:(bash zsh fish)'
    )

    _arguments -s $flags
}

_xlog "$@"

# Installation:
# Save to a directory in your $fpath, e.g., ~/.zsh/completions/_xlog
# Or add to ~/.zshrc:
# eval "$(xlog -completion zsh)"
`

const fishCompletionScript = `# fish completion for xlog

# Flag completions
complete -c xlog -l source -d 'Directory that will act as a storage' -r -F
complete -c xlog -l build -d 'Build all pages as static site in this directory' -r -F
complete -c xlog -l sitename -d 'Site name' -r
complete -c xlog -l index -d 'Index file name used as home page' -r
complete -c xlog -l notfoundpage -d 'Custom not found page' -r
complete -c xlog -l readonly -d 'Should xlog hide write operations'
complete -c xlog -l bind -d 'IP and port to bind the web server to' -r
complete -c xlog -l serve-insecure -d 'Accept http connections'
complete -c xlog -l csrf-cookie -d 'CSRF cookie name' -r
complete -c xlog -l disabled-extensions -d 'disable list of extensions' -r -a 'all'
complete -c xlog -l codestyle -d 'code highlighting style name' -r
complete -c xlog -l theme -d 'bulma theme to use' -r -a 'light dark'
complete -c xlog -l completion -d 'Generate shell completion script' -r -a 'bash zsh fish'

# Installation:
# Save to ~/.config/fish/completions/xlog.fish
# Or add to config.fish:
# xlog -completion fish | source
`
