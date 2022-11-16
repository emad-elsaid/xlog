Xlog is an HTTP server written in Go that serves markdown files as HTML and allows editing files online. 

![](/public/logo.png)

Xlog is a result of trying to build an offline personal knowledgebase with the ability to autolink pages together automatically. Without depending on proprietary file format or online service. 

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

# Getting started 

- Installation
- Custom installation
- Extensions
- Dependencies

# Principles

* Minimal design and dependencies
* Small core, flexible enough for developers to extend it.

# Documentation

- The README on Github will have basic usage and general information.
- This website serves as end user documentation and developer entry point for developing extensions
- There is also a Go package documentation that you can use to understand what xlog expose as public API

# License

Xlog is released under MIT license