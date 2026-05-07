# XLog

**Fast, Git-native framework for building knowledge bases**

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/xlog)](https://goreportcard.com/report/github.com/emad-elsaid/xlog) [![GoDoc](https://godoc.org/github.com/emad-elsaid/xlog?status.svg)](https://godoc.org/github.com/emad-elsaid/xlog)

XLog is a static site generator optimized for digital gardens and personal wikis. Written in Go with 37 built-in extensions for backlinks, hashtags, search, and more.

<p align="center"><img width="100%" src="public/screenshot.png" /></p>

## Features

- **Fast** - Written in Go, renders large knowledge bases in seconds
- **Live Preview** - Embedded web server with hot-reload shows changes instantly
- **Interconnected** - Automatic backlinks and bidirectional page relationships
- **Extensible** - 37 built-in extensions for photos, todos, search, and more
- **Git-Native** - Filesystem-based, works with any text editor and version control
- **Flexible Output** - Serve locally or generate static sites for deployment

## Quick Start

```bash
go install github.com/emad-elsaid/xlog/cmd/xlog@latest
mkdir my-notes
cd my-notes
echo "# Hello World" > index.md
xlog
# => Browse to http://localhost:3000
```

## Use Cases

- **Personal Wiki** - Interconnected notes with automatic backlinks
- **Research Notes** - Organize papers, citations, and ideas
- **Documentation** - Team knowledge bases and project docs
- **Learning Journal** - Study notes with hashtags and search
- **Digital Garden** - Public or private knowledge sharing

## How It Works

1. **Write** - Create markdown files in any text editor (Vim, Emacs, VS Code)
2. **Preview** - Run `xlog` to start the live preview server with hot-reload
3. **Enhance** - XLog automatically adds backlinks, hashtags, and search
4. **Deploy** - Generate static site with `xlog -build output/`

XLog runs a web server that watches your markdown files. When you click "Edit" in the browser, it opens the file in your configured editor. Save the file, and the browser automatically refreshes to show your changes.

## Why XLog?

**vs Other Static Generators:** XLog adds knowledge-base features (backlinks, hashtags, search) out of the box

**vs Cloud Tools:** XLog is free, self-hosted, privacy-focused, and works offline

**vs Obsidian Publish:** XLog is open source, customizable via extensions, and Git-native

**vs Notion:** XLog is local-first, markdown-based, and integrates with your existing workflow

## Documentation

- [Installation](https://xlog.emadelsaid.com/docs/Installation/)
- [Usage Guide](https://xlog.emadelsaid.com/docs/Usage/)
- [Extensions](https://xlog.emadelsaid.com/docs/extensions/)
- [Creating a Static Site](https://xlog.emadelsaid.com/tutorials/Creating%20a%20site)
- [Writing Extensions](https://xlog.emadelsaid.com/tutorials/Hello%20world%20extension/)
- [API Documentation](https://godoc.org/github.com/emad-elsaid/xlog)

## Extensions (37 Built-in)

XLog includes 37 extensions that enhance your knowledge base:

**Knowledge Base:**
- **Backlinks** - Automatic bidirectional links between pages
- **Hashtags** - Tag pages and browse by topic
- **Search** - Full-text search across all notes
- **Recent** - Activity tracking for recently modified pages

**Content:**
- **Photos** - EXIF data extraction for photo albums
- **Todos** - Task lists with checkboxes
- **Mermaid** - Diagrams and flowcharts
- **MathJax** - Mathematical notation

**And 29 more...** [View all extensions](https://xlog.emadelsaid.com/docs/extensions/)

## Shell Completion

```bash
# Bash
eval "$(xlog -completion bash)"

# Zsh
eval "$(xlog -completion zsh)"

# Fish
xlog -completion fish | source
```

## Contributing

- Create new extensions
- Improve the core codebase
- Report issues or suggest features

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## Vendored Packages

XLog vendors some dependencies for more control and stability. We thank the original authors:

*   **goldmark** by [Yusuke Inuzuka](https://github.com/yuin) - Core markdown parser
*   **goldmark-emoji** by [Yusuke Inuzuka](https://github.com/yuin) - Emoji support
*   **goldmark-highlighting** by [Yusuke Inuzuka](https://github.com/yuin) - Syntax highlighting
*   **goldmark-meta** by [Yusuke Inuzuka](https://github.com/yuin) - Frontmatter parsing
*   **goldmark-toc** by [Abhinav](https://github.com/abhinav) - Table of contents

## License

XLog is released under the [MIT license](LICENSE)
