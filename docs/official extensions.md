# XLog Extensions

XLog includes 38 built-in extensions that add features for knowledge bases, content creation, and developer workflows. All extensions are enabled by default (see [Usage](Usage.md) to disable specific ones).

## Knowledge Base Extensions

These extensions are specifically designed for building interconnected knowledge bases:

| Extension | Description |
|-----------|-------------|
| **Backlinks** (Autolink pages) | **Automatically converts page name mentions to links**. Creates bidirectional links between pages without manual linking. Core feature for networked knowledge. |
| **Hashtags** | **Add support for #hashtag syntax**. Organize notes by topic, browse all pages with a specific tag. Essential for categorization. |
| **Search** | **Full-text search across all pages**. Find content instantly without browsing. |
| **Recent** | **Lists pages ordered by last modified date**. Track recent activity and writing progress. |
| **Star** | **Pin important pages to footer**. Keep frequently accessed pages readily available. |
| **Date** | **Detects dates and links to timeline pages**. See all pages mentioning a specific date. |
| **Embed** | **Shortcode to embed one page inside another**. Transclude content across pages. |

## Content Creation Extensions

Enhance your writing with rich media and formatting:

| Extension | Description |
|-----------|-------------|
| **Photos** | Lists images in a directory Instagram-style. Great for photo journals. |
| **Images** | Display consecutive images in columns side-by-side instead of stacked. |
| **Upload file** | Upload files, screenshots, audio recordings, and camera photos directly. |
| **Link preview** | Preview tweets, Facebook posts, YouTube videos, Giphy links inline. |
| **MathJax** | Support mathematical notation with $ (inline) and $$ (blocks). |
| **Mermaid** | Create diagrams and flowcharts using MermaidJS syntax. |
| **TOC** | Automatic table of contents for long pages. |

## Writing & Editing Extensions

Tools for the writing workflow:

| Extension | Description |
|-----------|-------------|
| **Editor** | **Opens pages in your desktop text editor** (Vim, VS Code, Emacs, etc.). Core editing feature. |
| **Hotreload** | **Automatically reloads page when file changes on disk**. Instant preview of edits. |
| **Frontmatter** | YAML frontmatter support for metadata (title, tags, properties). |
| **Todo** | **Toggle checkboxes while viewing without edit mode**. Interactive task management. |
| **File operations** | Delete and rename pages from the UI. |
| **Heading** | Renders headings as linkable anchors. |

## Developer & Advanced Extensions

For technical documentation and advanced use cases:

| Extension | Description |
|-----------|-------------|
| **Blocks** | Define custom blocks using YAML code blocks. Extensible content types. |
| **Shortcode** | One-line and block shortcodes for custom functionality. |
| **sql_table** | Query long tables with SQL. Interactive data exploration. |
| **Github** | "Edit on GitHub" quick action for collaborative documentation. |
| **PGP** | Encrypt/decrypt .md.pgp files using GPG. Secure private notes. |

## Format Support Extensions

Support for additional file formats beyond markdown:

| Extension | Description |
|-----------|-------------|
| **HTML** | Treat HTML files as pages (.html, .htm, .xhtml). |
| **Pandoc** | Render .org, .rst, .rtf, .odt files using Pandoc. |

## Publishing & Integration Extensions

Deploy and integrate with external services:

| Extension | Description |
|-----------|-------------|
| **RSS** | RSS feed at `/+/feed.rss`, added to page headers. |
| **Sitemap** | sitemap.xml for search engine crawling. |
| **Opengraph** | Open Graph meta tags for social media previews. |
| **Manifest** | Web app manifest.json for PWA support. |
| **ActivityPub** | Webfinger and ActivityPub actor for Fediverse integration. |
| **Disqus** | Add Disqus comments to pages. |
| **Giscus** | Add Giscus comments powered by GitHub Discussions. |

## Formatting Extensions

Improve text rendering and display:

| Extension | Description |
|-----------|-------------|
| **Autolink** | Shorten long link strings to save space. |
| **RTL** | Proper text direction for right-to-left languages. |

## Custom Content Extensions

Add arbitrary HTML to pages:

| Extension | Description |
|-----------|-------------|
| **Custom Widget** | Inject custom HTML in `<head>`, before, or after page content. |

---

## Extension Showcase

### Backlinks in Action

When you mention `[[Another Page]]` in your markdown:

```markdown
This concept relates to [[Graph Theory]] and [[Network Effects]].
```

XLog automatically:
1. Creates clickable links to those pages
2. Shows on "Graph Theory" that this page links to it (backlink)
3. Builds a knowledge graph of relationships

### Hashtags for Organization

Use hashtags to categorize:

```markdown
# My Research Notes

This is about #machine-learning and #neural-networks.

Some thoughts on #attention-mechanisms in transformers.
```

Click any hashtag to see all pages tagged with it.

### Search Everything

Press Ctrl+K or click Search to find any content across your knowledge base instantly.

### Live Preview Workflow

1. Click "Edit" on any page
2. Opens in your configured editor (Vim, VS Code, etc.)
3. Make changes and save
4. Browser automatically reloads with updated content
5. No manual refresh needed

### Giscus Comments

Enable GitHub Discussions-powered comments on your pages:

```bash
xlog -giscus-repo "owner/repo" \
     -giscus-repo-id "R_kgDOG1234" \
     -giscus-category "Announcements" \
     -giscus-category-id "DIC_kwDOG5678" \
     -giscus-mapping "pathname" \
     -giscus-theme "preferred_color_scheme" \
     -giscus-lang "en"
```

To get your repository and category IDs, visit [giscus.app](https://giscus.app) and follow the configuration wizard. You'll need:

1. A public GitHub repository with Discussions enabled
2. The [giscus app](https://github.com/apps/giscus) installed on that repository
3. A Discussion category (recommended: Announcements type)

Optional parameters:
- `-giscus-mapping`: How to match pages to discussions (pathname, url, title, og:title)
- `-giscus-theme`: Visual theme (light, dark, preferred_color_scheme, etc.)
- `-giscus-lang`: Interface language (en, ar, zh, etc.)

## Disabling Extensions

To disable specific extensions:

```bash
xlog -disabled-extensions "photos,videos,activitypub"
```

Comma-separated list of extension names (case-sensitive, use lowercase).

## Extension Architecture

Extensions use defined extension points to add functionality. See [extensions.md](extensions.md) for developer documentation on creating custom extensions.

All built-in extensions are defined under `/extensions` and automatically imported in `cmd/xlog/xlog.go`.

## Most Valuable Extensions for Knowledge Bases

If you're building a personal wiki or digital garden, these extensions provide the core experience:

1. **Backlinks** - Automatic bidirectional linking
2. **Hashtags** - Topic organization
3. **Search** - Find anything instantly
4. **Recent** - Track writing activity
5. **Editor** - Edit in your preferred tool
6. **Hotreload** - Instant preview
7. **Star** - Pin important pages

Together, these create a powerful knowledge base system that rivals dedicated tools like Roam Research or Obsidian, while keeping all your data in portable markdown files.

## See Also

- [Creating Extensions](extensions.md) - Developer guide for building custom extensions
- [Usage](Usage.md) - How to enable/disable extensions
- [Why XLog](Why-XLog.md) - Understanding XLog's knowledge base focus
