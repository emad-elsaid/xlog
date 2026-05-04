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

    # All flags from xlog -h (complete list)
    opts="-activitypub.domain -activitypub.icon -activitypub.image -activitypub.summary -activitypub.username"
    opts="$opts -bind -build -codestyle -completion -csrf-cookie"
    opts="$opts -custom.after_view -custom.before_view -custom.head"
    opts="$opts -disabled-extensions -disqus -editor"
    opts="$opts -github.url -gpg -html -index -notfoundpage"
    opts="$opts -og.domain -pandoc -readonly"
    opts="$opts -rss.description -rss.domain -rss.limit"
    opts="$opts -serve-insecure -sitemap.domain -sitename -source"
    opts="$opts -sql-table.threshold -theme -twitter.username"

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
        -custom.after_view|-custom.before_view|-custom.head|-activitypub.icon|-activitypub.image)
            COMPREPLY=( $(compgen -f -- ${cur}) )
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
        '-activitypub.domain[Domain used for activitypub stream absolute URLs]:domain:'
        '-activitypub.icon[Path to the activitypub profile icon]:file:_files'
        '-activitypub.image[Path to the activitypub profile image]:file:_files'
        '-activitypub.summary[Summary of the user for activitypub actor]:summary:'
        '-activitypub.username[Username for activitypub actor]:username:'
        '-bind[IP and port to bind the web server to]:address:'
        '-build[Build all pages as static site in this directory]:directory:_directories'
        '-codestyle[Code highlighting style name]:style:'
        '-completion[Generate shell completion script]:shell:(bash zsh fish)'
        '-csrf-cookie[CSRF cookie name]:name:'
        '-custom.after_view[Path to file included AFTER page content]:file:_files'
        '-custom.before_view[Path to file included BEFORE page content]:file:_files'
        '-custom.head[Path to file included in every page <head> tag]:file:_files'
        '-disabled-extensions[Disable list of extensions by name, comma separated]:extensions:(all)'
        '-disqus[Disqus domain name]:domain:'
        '-editor[Command to use to open pages for editing]:command:'
        '-github.url[Repository URL for edit on GitHub quick action]:url:'
        '-gpg[PGP key ID to decrypt and edit .md.pgp files]:key-id:'
        '-html[Consider HTML files as pages]'
        '-index[Index file name used as home page]:filename:'
        '-notfoundpage[Custom not found page]:filename:'
        '-og.domain[Domain for OpenGraph meta tags]:domain:'
        '-pandoc[Path to pandoc binary for converting various formats]:path:_files'
        '-readonly[Hide write operations]'
        '-rss.description[RSS feed description]:description:'
        '-rss.domain[Domain for RSS feed absolute URLs]:domain:'
        '-rss.limit[Maximum number of items in RSS feed]:limit:'
        '-serve-insecure[Accept HTTP connections]'
        '-sitemap.domain[Domain for sitemap XML absolute URLs]:domain:'
        '-sitename[Site name]:name:'
        '-source[Directory that will act as storage]:directory:_directories'
        '-sql-table.threshold[Minimum rows to render SQL table]:threshold:'
        '-theme[Bulma theme to use]:theme:(light dark)'
        '-twitter.username[Twitter username for Twitter card meta tags]:username:'
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

# ActivityPub flags
complete -c xlog -l activitypub.domain -d 'Domain used for activitypub stream absolute URLs' -r
complete -c xlog -l activitypub.icon -d 'Path to the activitypub profile icon' -r -F
complete -c xlog -l activitypub.image -d 'Path to the activitypub profile image' -r -F
complete -c xlog -l activitypub.summary -d 'Summary of the user for activitypub actor' -r
complete -c xlog -l activitypub.username -d 'Username for activitypub actor' -r

# Server and build flags
complete -c xlog -l bind -d 'IP and port to bind the web server to' -r
complete -c xlog -l build -d 'Build all pages as static site in this directory' -r -F
complete -c xlog -l codestyle -d 'Code highlighting style name' -r
complete -c xlog -l completion -d 'Generate shell completion script' -r -a 'bash zsh fish'
complete -c xlog -l csrf-cookie -d 'CSRF cookie name' -r

# Custom content flags
complete -c xlog -l custom.after_view -d 'Path to file included AFTER page content' -r -F
complete -c xlog -l custom.before_view -d 'Path to file included BEFORE page content' -r -F
complete -c xlog -l custom.head -d 'Path to file included in every page <head> tag' -r -F

# Extension and integration flags
complete -c xlog -l disabled-extensions -d 'Disable list of extensions by name, comma separated' -r -a 'all'
complete -c xlog -l disqus -d 'Disqus domain name' -r
complete -c xlog -l editor -d 'Command to use to open pages for editing' -r

# GitHub and encryption flags
complete -c xlog -l github.url -d 'Repository URL for edit on GitHub quick action' -r
complete -c xlog -l gpg -d 'PGP key ID to decrypt and edit .md.pgp files' -r

# Content type flags
complete -c xlog -l html -d 'Consider HTML files as pages'

# Page configuration flags
complete -c xlog -l index -d 'Index file name used as home page' -r
complete -c xlog -l notfoundpage -d 'Custom not found page' -r

# SEO and metadata flags
complete -c xlog -l og.domain -d 'Domain for OpenGraph meta tags' -r
complete -c xlog -l pandoc -d 'Path to pandoc binary for converting various formats' -r -F
complete -c xlog -l readonly -d 'Hide write operations'

# RSS flags
complete -c xlog -l rss.description -d 'RSS feed description' -r
complete -c xlog -l rss.domain -d 'Domain for RSS feed absolute URLs' -r
complete -c xlog -l rss.limit -d 'Maximum number of items in RSS feed' -r

# Server mode and sitemap flags
complete -c xlog -l serve-insecure -d 'Accept HTTP connections'
complete -c xlog -l sitemap.domain -d 'Domain for sitemap XML absolute URLs' -r
complete -c xlog -l sitename -d 'Site name' -r
complete -c xlog -l source -d 'Directory that will act as storage' -r -F

# Extension-specific flags
complete -c xlog -l sql-table.threshold -d 'Minimum rows to render SQL table' -r

# Theme and social flags
complete -c xlog -l theme -d 'Bulma theme to use' -r -a 'light dark'
complete -c xlog -l twitter.username -d 'Twitter username for Twitter card meta tags' -r

# Installation:
# Save to ~/.config/fish/completions/xlog.fish
# Or add to config.fish:
# xlog -completion fish | source
`
