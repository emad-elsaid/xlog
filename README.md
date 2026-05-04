XLog
=========

:vhs: Xlog is a static site generator for digital gardening written in Go. It serves markdown files as HTML and allows editing files online. It focuses on enriching markdown files and surfacing implicit links between pages.

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/xlog)](https://goreportcard.com/report/github.com/emad-elsaid/xlog) [![GoDoc](https://godoc.org/github.com/emad-elsaid/xlog?status.svg)](https://godoc.org/github.com/emad-elsaid/xlog)

<p align="center"><img width="200" src="public/logo.png" /></p>

<p align="center"><img width="100%" src="public/screenshot.png" /></p>


# Documentation

* [Documentation](https://xlog.emadelsaid.com/)
* [Installation](https://xlog.emadelsaid.com/docs/Installation/)
* [Usage](https://xlog.emadelsaid.com/docs/Usage/)
* [Generating static site](https://xlog.emadelsaid.com/tutorials/Creating%20a%20site)
* [Overriding Assets](https://xlog.emadelsaid.com/docs/Assets)
* [Extensions](https://xlog.emadelsaid.com/docs/extensions/)
* [Writing Your Own Extension](https://xlog.emadelsaid.com/tutorials/Hello%20world%20extension/)

# Shell Completion

Xlog supports shell completion for bash, zsh, and fish. This enables tab completion for all command-line flags.

**Bash:**
```bash
# Add to ~/.bashrc
eval "$(xlog -completion bash)"
```

**Zsh:**
```bash
# Add to ~/.zshrc
eval "$(xlog -completion zsh)"
```

**Fish:**
```fish
# Add to ~/.config/fish/config.fish
xlog -completion fish | source
```

# Vendored Packages

Xlog vendors some of its dependencies for more control over the changes and to allow for major refactoring needed for the project. We would like to thank the original authors for their great work.

*   **goldmark** by [Yusuke Inuzuka](https://github.com/yuin): The core markdown parser. Vendored from [github.com/yuin/goldmark](https://github.com/yuin/goldmark).
*   **goldmark-emoji** by [Yusuke Inuzuka](https://github.com/yuin): For adding emoji support. Vendored from [github.com/yuin/goldmark-emoji](https://github.com/yuin/goldmark-emoji).
*   **goldmark-highlighting** by [Yusuke Inuzuka](https://github.com/yuin): For syntax highlighting. Vendored from [github.com/yuin/goldmark-highlighting](https://github.com/yuin/goldmark-highlighting).
*   **goldmark-meta** by [Yusuke Inuzuka](https://github.com/yuin): For frontmatter parsing. Vendored from [github.com/yuin/goldmark-meta](https://github.com/yuin/goldmark-meta).
*   **goldmark-toc** by [Abhinav](https://github.com/abhinav): For generating a table of contents. Vendored from [github.com/abhinav/goldmark-toc](https://github.com/abhinav/goldmark-toc).

# License

Xlog is released under [MIT license](LICENSE)
