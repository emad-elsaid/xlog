```hero
image: /public/logo.png
title: Fast, Git-native framework for building knowledge bases
subtitle: XLog is a static site generator optimized for digital gardens and personal wikis. Written in Go with 37 built-in extensions for backlinks, hashtags, search, and more.
```

# ⚡ Quick Start

```shell
go install github.com/emad-elsaid/xlog/cmd/xlog@latest
mkdir my-notes
cd my-notes
echo "# Hello World" > index.md
xlog
# => Browse to http://localhost:3000
```

# 🚀 Why XLog?

## Optimized for Knowledge Work

Written in Go, optimized for markdown-heavy sites with bidirectional links. With its live preview server and hot-reload, XLog lets you write in your favorite editor and see changes instantly.

## Flexible Framework

With its 37 built-in extensions (hashtags, backlinks, search, photos, todos), XLog is widely used for personal wikis, research notes, documentation, and learning journals.

## Developer-Friendly

Extension API for customization, filesystem-based (no database), Git-native workflow, and works with any text editor (Vim, Emacs, VS Code).

## Embedded Web Server

Use XLog's embedded web server during writing to instantly see changes. Generate static sites for deployment to GitHub Pages, Netlify, or any host.

---

# 🔑 Core Features

- **Fast** - Renders large knowledge bases in seconds
- **Live Preview** - Hot-reload shows changes as you save
- **Interconnected** - Automatic backlinks and page relationships
- **Extensible** - 37 built-in extensions, easy to add more
- **Git-Native** - Filesystem-based, version control friendly
- **Flexible Output** - Serve locally or generate static sites

# 📊 How XLog Compares

| Feature | XLog | Other Static Generators | Cloud Tools |
|---------|------|------------------------|-------------|
| Backlinks | ✅ Automatic | ❌ Manual | ✅ Built-in |
| Live Preview | ✅ Hot-reload | ✅ Dev server | ✅ Cloud |
| Editor | Desktop (Any) | Desktop (Any) | Browser |
| Hosting | Self/Static | Static | Cloud only |
| Price | Free | Free | $8-16/mo |
| Privacy | 100% Local | 100% Local | Cloud |

**XLog sits between pure static generators and cloud-based tools** - offering knowledge-base features with complete ownership of your data.

# 📌 Use Cases

- **Personal Wiki** - Interconnected notes with automatic backlinks
- **Research Notes** - Organize papers, citations, and ideas  
- **Documentation** - Team knowledge bases and project docs
- **Learning Journal** - Study notes with hashtags and search
- **Digital Garden** - Public or private knowledge sharing

# 🎯 How It Works

1. **Write** - Create markdown files in your text editor
2. **Preview** - Run `xlog` for live preview with hot-reload
3. **Enhance** - Automatic backlinks, hashtags, and search
4. **Deploy** - Generate static site with `xlog -build output/`

**Editor Workflow:** When you click "Edit" in the browser, XLog opens the file in your configured editor (VS Code, Vim, Emacs, etc). Save the file, and the browser automatically refreshes.

# ⚖️ Principles

* Filesystem-based - No databases required
* Minimal design and dependencies
* Small core, flexible extension system
* Enhance existing markdown syntax, don't invent new syntax

# 🌱 Getting Started

- Installation
- Usage
- Creating a Static Site
- Extensions

# 🧩 Extensions (37 Built-in)

/hashtag-pages-grid extension

# 🎓 Tutorials

/hashtag-pages-grid tutorial

# 🤝 Contributors

/github-user name: juanolon
/github-user name: m4salah
/github-user name: ProvoK
/github-user name: scroot
/github-user name: disconnect3d
/github-user name: mohamed-zezo
/github-user name: aradwann
/github-user name: aaelsay3d

## Contributing

- Create new extensions
- Improve the core codebase
- Report issues or suggest features

# 🧑‍🤝‍🧑 Community

- :left_speech_bubble: [Discussions](https://github.com/emad-elsaid/xlog/discussions)
- :beetle: [Issues](https://github.com/emad-elsaid/xlog/issues)
- :keyboard: [GitHub](https://github.com/emad-elsaid/xlog)

# 📜 License

XLog is released under MIT license
