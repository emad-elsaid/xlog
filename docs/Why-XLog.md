# Why Choose XLog?

XLog is a fast, Git-native framework for building knowledge bases. This guide helps you decide if XLog is the right tool for your needs.

## When to Choose XLog

### You Want a Personal Knowledge Base

XLog excels at personal wikis, research notes, and digital gardens. With automatic backlinks and hashtags, you can build an interconnected knowledge base where ideas naturally link together.

**Perfect for:**
- Personal wikis and documentation
- Research notes with bidirectional links
- Learning journals and study notes
- Digital gardens for knowledge cultivation
- Zettelkasten-style note systems

### You Prefer Desktop Editors

XLog works with your favorite text editor (Vim, Emacs, VS Code, Sublime). Edit markdown files locally, and XLog's live preview server shows changes instantly. No browser-based editor means you use the tools you already know.

**Perfect for:**
- Developers comfortable with text editors
- Writers who prefer markdown workflows
- Anyone who wants editor choice and flexibility

### You Value Git-Native Workflow

XLog stores everything as markdown files in your filesystem. This means full Git integration, version control, and the ability to sync across devices using standard Git workflows.

**Perfect for:**
- Version-controlled documentation
- Collaborative wikis using Git
- Backing up notes to GitHub/GitLab
- Offline-first knowledge management

### You Need Fast Performance

Written in Go with minimal dependencies, XLog renders large knowledge bases in seconds. The embedded web server provides instant hot-reload during writing.

**Perfect for:**
- Large knowledge bases (thousands of pages)
- Quick iteration and writing flow
- Low-resource environments

## When NOT to Choose XLog

### You Want Browser-Based Editing

XLog requires a desktop text editor. If you need to edit from any browser without installing software, consider cloud-based tools like Notion or Obsidian Publish.

**Consider instead:**
- Notion (cloud-based, browser editing)
- Obsidian Publish (browser interface option)
- Wiki.js (self-hosted with browser editor)

### You Need Real-Time Collaboration

XLog uses Git for collaboration, not real-time co-editing. If multiple people need to edit simultaneously, you'll need a different solution.

**Consider instead:**
- Notion (real-time collaboration)
- Google Docs (simultaneous editing)
- HedgeDoc (collaborative markdown)

### You Want WYSIWYG Editing

XLog works with markdown source files. If you or your team prefer visual editing without markdown syntax, XLog isn't the right fit.

**Consider instead:**
- Notion (block-based editor)
- Obsidian (has live preview mode)
- Typora (WYSIWYG markdown)

### You Need a CMS for Public Websites

While XLog can generate static sites, it's optimized for knowledge bases, not general websites or blogs. If you need complex layouts, themes, or CMS features, use a dedicated static site generator.

**Consider instead:**
- Hugo (general-purpose static sites)
- Jekyll (blogs and documentation)
- Eleventy (flexible static sites)

## XLog's Philosophy

XLog is built on three core principles:

### 1. Filesystem-Based

Everything is markdown files in folders. No database, no lock-in. Your notes are portable, readable, and future-proof.

### 2. Git-Native

Version control is built into the workflow. Every change is trackable, recoverable, and syncable using standard Git tools.

### 3. Editor-Agnostic

Use any text editor you want. XLog doesn't force you into a specific editing environment—it works with the tools you already know.

## XLog's Unique Value

What makes XLog different from other static site generators and knowledge base tools:

### Automatic Backlinks

Unlike most static generators, XLog automatically detects when pages link to each other and shows bidirectional relationships. This creates a networked knowledge graph without manual effort.

### Live Preview + Static Generation

Get both worlds: fast development with live preview during writing, and static HTML for deployment. No build step needed during editing.

### 37 Built-In Extensions

Knowledge base features out of the box: backlinks, hashtags, search, todos, photos, and more. No plugin hunting required.

### Lightweight & Fast

Single binary, minimal dependencies, fast rendering. Install in seconds, build in seconds.

## Target Use Cases

XLog is specifically designed for:

### Personal Wikis
Interconnected personal knowledge bases with automatic backlinks and hashtags for organization.

### Research Notes
Academic research, literature reviews, and study notes with bidirectional linking between concepts.

### Documentation Sites
Technical documentation, API docs, and internal wikis with Git-based version control.

### Digital Gardens
Cultivate ideas over time in an interconnected web of thoughts, perfect for learning in public.

### Learning Journals
Track learning progress, connect concepts, and build knowledge incrementally.

## Quick Comparison

| Aspect | XLog | Obsidian | Notion | Hugo |
|--------|------|----------|--------|------|
| **Editing** | Desktop editor | Desktop app | Browser | Desktop editor |
| **Storage** | Local files | Local files | Cloud database | Local files |
| **Backlinks** | Automatic | Automatic | Manual | Manual |
| **Collaboration** | Git-based | Sync plugins | Real-time | Git-based |
| **Deployment** | Static HTML | Obsidian Publish | Cloud-only | Static HTML |
| **Cost** | Free, open-source | Free (Publish paid) | Free tier, then paid | Free, open-source |

## Getting Started

If XLog sounds like the right fit, check out:

- [Installation](Installation.md) - Install XLog in under 2 minutes
- [Workflow](Workflow.md) - Understand the editor + preview workflow
- [Comparison](Comparison.md) - Detailed comparisons with alternatives

## Still Unsure?

Ask yourself:
- Do I prefer markdown files over proprietary formats? → **Yes = XLog**
- Do I want Git version control for my notes? → **Yes = XLog**
- Do I need automatic backlinks and hashtags? → **Yes = XLog**
- Do I want to use my own text editor? → **Yes = XLog**
- Do I need real-time browser collaboration? → **No = Consider alternatives**

XLog is for people who value local-first, Git-native, editor-agnostic knowledge base building with automatic backlinks. If that's you, welcome!
