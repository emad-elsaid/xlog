# XLog Comparison Guide

When choosing a tool for building your knowledge base or digital garden, understanding the tradeoffs is crucial. This guide compares XLog with other popular options to help you make an informed decision.

---

## XLog vs Other Static Site Generators

### When to Choose XLog

**Choose XLog if you want:**
- Automatic backlinks between pages
- Hashtag-based organization out of the box
- Full-text search without plugins
- Live preview with instant hot-reload
- Knowledge-base features ready to use

**Choose other generators if you want:**
- General-purpose website building (not knowledge-focused)
- Massive theme ecosystem
- Hosted CMS admin panel
- Mature plugin marketplace

### Feature Comparison

| Feature | XLog | Other Go Generators | Ruby Generators | JS Generators |
|---------|------|---------------------|-----------------|---------------|
| **Speed** | Fast (Go) | Fast (Go) | Slow (Ruby) | Medium (Node) |
| **Backlinks** | ✅ Automatic | ❌ Manual | ❌ Manual | ❌ Manual |
| **Hashtags** | ✅ Built-in | ⚠️ Plugin | ⚠️ Plugin | ⚠️ Plugin |
| **Search** | ✅ Built-in | ⚠️ External | ⚠️ External | ⚠️ External |
| **Live Preview** | ✅ Hot-reload | ✅ Dev server | ✅ Dev server | ✅ Dev server |
| **Extensions** | 37 built-in | Themes/plugins | Plugins | Plugins |
| **Setup** | Single binary | Single binary | Ruby + gems | Node + npm |
| **Focus** | Knowledge bases | General sites | Blogs | General sites |

### Use Case Guide

**Personal Wiki / Research Notes**
- ✅ **XLog** - Backlinks and hashtags make this ideal
- ⚠️ Other generators - Require manual linking

**Blog / Marketing Site**
- ⚠️ XLog - Can work but not optimized for this
- ✅ Other generators - Better theme options

**Documentation Site**
- ✅ **XLog** - Search and backlinks help navigation
- ✅ Other generators - Also good choice

**Portfolio / Landing Page**
- ❌ XLog - Too knowledge-base focused
- ✅ Other generators - Better suited

---

## XLog vs Cloud-Based Tools

### XLog vs Obsidian Publish

| Feature | XLog | Obsidian Publish |
|---------|------|------------------|
| **Price** | Free | $8-16/month |
| **Hosting** | Self-hosted | Cloud |
| **Editor** | Any (Vim, VS Code, etc) | Obsidian app |
| **Backlinks** | ✅ Automatic | ✅ Automatic |
| **Graph View** | ⚠️ Basic | ✅ Advanced |
| **Offline** | ✅ 100% | ✅ With app |
| **Privacy** | ✅ 100% local | Cloud-based |
| **Customization** | ✅ 37 extensions | ⚠️ Limited |
| **Deployment** | Git/GitHub Pages | Automatic |

**Choose XLog if:**
- You want free, self-hosted solution
- You prefer your own text editor
- You want complete control over hosting
- Privacy is critical

**Choose Obsidian Publish if:**
- You're already using Obsidian
- You want zero-config publishing
- Advanced graph view is important
- You don't mind monthly fee

### XLog vs Notion

| Feature | XLog | Notion |
|---------|------|--------|
| **Price** | Free | Free-$10/month |
| **Editor** | Desktop (any) | Browser |
| **Format** | Markdown files | Proprietary DB |
| **Offline** | ✅ Full | ⚠️ Limited |
| **Privacy** | ✅ 100% local | Cloud-based |
| **Speed** | Fast | Medium |
| **Git Integration** | ✅ Native | ❌ Export only |
| **Export** | ✅ Files already yours | ⚠️ Manual export |
| **Collaboration** | Git-based | ✅ Real-time |

**Choose XLog if:**
- You want local-first, markdown-based
- Git workflow is important
- You prefer text editor over browser
- Long-term data ownership matters

**Choose Notion if:**
- Real-time collaboration needed
- You prefer WYSIWYG editing
- Databases and kanban boards required
- Browser-based workflow is fine

---

## Migration Paths

### From Other Generators to XLog

**If your current generator has:**
- Markdown files → XLog reads them directly
- Frontmatter → XLog supports YAML frontmatter
- Images/assets → Copy to XLog's public folder
- Custom themes → Adapt templates (see docs)

**Migration checklist:**
1. Copy markdown files to XLog directory
2. Move images to `public/` folder
3. Update image paths if needed
4. Run `xlog` to preview
5. Adjust frontmatter if needed

### From Obsidian to XLog

1. Point XLog to your Obsidian vault: `xlog -source ~/Obsidian/MyVault`
2. XLog reads wikilinks and backlinks automatically
3. Attachments work if in vault folder
4. Use both tools simultaneously (Obsidian for editing, XLog for publishing)

### From Notion to XLog

1. Export Notion workspace as Markdown & CSV
2. Copy markdown files to XLog directory
3. Move images to `public/` folder
4. Update image references
5. Convert Notion databases to markdown tables

---

## Decision Matrix

### Quick Selector

**I want to...**

**"Build a personal wiki with backlinks"**
→ ✅ XLog | ✅ Obsidian Publish | ⚠️ Other generators (need plugins)

**"Write technical documentation"**
→ ✅ XLog | ✅ Other generators (both good)

**"Create a blog or marketing site"**
→ ⚠️ XLog (can work) | ✅ Other generators (better themes)

**"Organize research notes"**
→ ✅ XLog | ✅ Notion | ✅ Obsidian

**"Work 100% offline"**
→ ✅ XLog | ✅ Other generators | ❌ Cloud tools

**"Collaborate in real-time"**
→ ❌ XLog (Git-based) | ✅ Notion | ✅ Cloud tools

**"Pay nothing, ever"**
→ ✅ XLog | ✅ Other generators | ⚠️ Cloud tools (free tiers limited)

**"Use my favorite text editor"**
→ ✅ XLog | ✅ Other generators | ❌ Cloud tools (browser-based)

---

## Summary

**XLog is best for:**
- Developers building personal wikis
- Researchers organizing interconnected notes
- Privacy-conscious users wanting local-first
- Git-native workflows
- Knowledge bases needing backlinks + search out of the box

**XLog is NOT ideal for:**
- Non-technical users unfamiliar with terminal
- Real-time team collaboration needs
- Marketing sites needing fancy themes
- Users wanting WYSIWYG browser editing

**Bottom line:** XLog sits between general-purpose static generators (which lack knowledge-base features) and cloud tools (which lack privacy/control). If you value local-first, Git-native, and automatic backlinks, XLog is the right choice.

---

## Getting Help

If you're still unsure, ask in our community:
- [GitHub Discussions](https://github.com/emad-elsaid/xlog/discussions)
- [Issues](https://github.com/emad-elsaid/xlog/issues)

We're happy to help you decide if XLog fits your needs!
