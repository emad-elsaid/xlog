:vhs: Xlog is a static site generator for digital gardening written in Go. It serves markdown files as HTML and allows editing files online. It focuses on enriching markdown files and surfacing implicit links between pages.

![](/public/logo.png)

Xlog is a result of trying to build an offline personal knowledgebase with the ability to autolink pages together automatically. Without depending on proprietary file format or online service. 

# :runner: Quick Start

```shell
go install github.com/emad-elsaid/xlog/cmd/xlog@latest
mkdir new-site
cd new-site
xlog
# => Now browse to http://localhost:3000
```

# Core Features

- Serves any file from current directory
- Any markdown is rendered to HTML format
- Supports Github flavor markdown
- Has a list of tools defined by extensions. triggered with `Ctrl+K`
- Use image at the top of the page as a cover image
- The first Emoji used in the page will be considered the icon of the page and displayed beside the title
- Shows task list (Done/Total tasks) beside page link

# Usecases

- Local server for Note taking or digital gardening
- Generate static website just like the one you're reading now

# :checkered_flag: Getting started 

- Installation
- Custom installation
- Usage
- Creating a site
- Extensions
- Dependencies
- Security

# :scroll: Principles

* Uses the file system. No databases required
* Minimal design and dependencies
* Small core, flexible enough for developers to extend it.
* Avoid adding syntax to markdown, instead enhance how existing syntax is rendered

# :book: Documentation

- This website serves as end user documentation and developer entry point for developing extensions
- There is also a Go package documentation that you can use to understand what xlog expose as public API

# :bulb: Tutorials

- Hello world extension
- Create your own digital garden on Github

# Contributing

You can help Xlog in many ways:

- Create a new extension
- Improve the core codebase
- Package it for different operating systems or different Linux distribution

# :people_holding_hands: Community

- :left_speech_bubble: [Discussing ideas](https://github.com/emad-elsaid/xlog/discussions)
- :beetle: [Reporting issues](https://github.com/emad-elsaid/xlog/issues)
- :keyboard: [Contributors](https://github.com/emad-elsaid/xlog/graphs/contributors)
- Github

# License

Xlog is released under MIT license