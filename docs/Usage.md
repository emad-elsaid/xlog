# Usage

XLog provides a live preview server for writing and a static site generator for deployment. This guide covers common usage patterns and command-line options.

## Basic Usage

### Start the Live Preview Server

The most common way to use XLog is with the live preview server:

```bash
cd my-notes
xlog
# => Server running at http://localhost:3000
```

This starts a web server that:
- Displays your markdown files as web pages
- Provides hot-reload (changes appear instantly when you save)
- Lets you click "Edit" to open files in your desktop editor
- Shows automatic backlinks and page relationships

```callout-info
**Desktop Editing Workflow**  
XLog doesn't provide browser-based editing. Instead, clicking "Edit" opens the file in your configured text editor (Vim, VS Code, Emacs, etc.). You edit locally, save the file, and see changes instantly in your browser.
```

### Generate a Static Site

To deploy your knowledge base as static HTML:

```bash
xlog -build ./output
# => Static site generated in ./output directory
```

This creates a complete static website you can deploy to:
- GitHub Pages
- Netlify
- Vercel
- Any static hosting service

## Command-Line Flags

### Essential Flags

#### `-editor` - Set Your Text Editor

Specify which editor opens when you click "Edit":

```bash
# Use VS Code
xlog -editor "code"

# Use Vim
xlog -editor "vim"

# Use Emacs
xlog -editor "emacs"

# Use Sublime Text
xlog -editor "subl -w"
```

If not specified, XLog uses the `$EDITOR` environment variable. Set it in your shell config:

```bash
export EDITOR="vim"  # Add to ~/.bashrc or ~/.zshrc
```

**Examples for different editors:**

| Editor | Command |
|--------|---------|
| VS Code | `code` or `code --wait` |
| Vim | `vim` |
| Emacs | `emacs` or `emacsclient -n -a emacs` |
| Sublime Text | `subl -w` |
| Nano | `nano` |
| Notepad++ (Windows) | `notepad++` |

#### `-bind` - Change Server Address

By default, XLog binds to `127.0.0.1:3000`. Change it with:

```bash
# Use port 8080
xlog -bind :8080

# Bind to all interfaces (accessible from network)
xlog -bind 0.0.0.0:3000
```

```callout-warning
**Security Note**  
Binding to `0.0.0.0` makes XLog accessible from your network. Only use this on trusted networks or with proper authentication.
```

#### `-source` - Specify Content Directory

By default, XLog uses the current directory. Change it with:

```bash
xlog -source ~/my-knowledge-base
```

#### `-build` - Generate Static Site

Create a static HTML version for deployment:

```bash
# Build to ./public directory
xlog -build ./public

# Build to custom location
xlog -build ~/sites/my-wiki
```

### Customization Flags

#### `-sitename` - Set Site Name

Appears in the header and title tags:

```bash
xlog -sitename "My Knowledge Base"
```

#### `-theme` - Choose Color Theme

```bash
# Light theme
xlog -theme light

# Dark theme
xlog -theme dark

# System preference (default)
xlog -theme ""
```

#### `-index` - Set Homepage File

Default is `index.md`. Change it with:

```bash
xlog -index "home"  # Uses home.md as homepage
```

#### `-notfoundpage` - Custom 404 Page

```bash
xlog -notfoundpage "custom-404"  # Uses custom-404.md for 404 errors
```

### Content Processing Flags

#### `-html` - Include HTML Files

By default, XLog only processes markdown (.md) files. Enable HTML:

```bash
xlog -html
```

#### `-pandoc` - Use Pandoc for Additional Formats

Render .org, .rst, .rtf, .odt files using Pandoc:

```bash
xlog -pandoc
```

Requires Pandoc to be installed.

#### `-codestyle` - Syntax Highlighting Theme

Change code block syntax highlighting:

```bash
xlog -codestyle monokai
xlog -codestyle github
xlog -codestyle dracula  # default
```

See [Chroma styles](https://pkg.go.dev/github.com/alecthomas/chroma/v2/styles) for available themes.

### Extension Management

#### `-disabled-extensions` - Disable Extensions

Disable specific extensions by name (comma-separated):

```bash
xlog -disabled-extensions "photos,videos,backlinks"
```

See [Extensions](extensions.md) for available extension names.

### Advanced Flags

#### `-readonly` - Read-Only Mode

Disable all write operations (editing, deleting, creating):

```bash
xlog -readonly
```

Useful for public deployments where you only want to display content.

#### `-serve-insecure` - Accept HTTP Connections

By default, XLog expects HTTPS for CSRF protection. Allow HTTP:

```bash
xlog -serve-insecure
```

```callout-warning
Only use this for local development. Production deployments should use HTTPS.
```

#### `-gpg` - Encrypt Pages with GPG

Encrypt/decrypt .md.pgp files using GPG:

```bash
xlog -gpg "YOUR_KEY_ID"
```

#### `-custom.head` - Inject Custom HTML

Include custom HTML in the `<head>` tag of every page:

```bash
xlog -custom.head ./custom-head.html
```

#### `-custom.before_view` - Content Before Page

Include custom HTML before page content:

```bash
xlog -custom.before_view ./header.html
```

#### `-custom.after_view` - Content After Page

Include custom HTML after page content:

```bash
xlog -custom.after_view ./footer.html
```

### Integration Flags

#### GitHub Integration

Enable "Edit on GitHub" buttons:

```bash
xlog -github.url "https://github.com/username/repo/edit/master/docs"
```

#### RSS Feed Configuration

```bash
# Set RSS domain (without https://)
xlog -rss.domain "example.com"

# Set RSS description
xlog -rss.description "My knowledge base feed"

# Limit feed items
xlog -rss.limit 50  # default: 30
```

#### Open Graph Meta Tags

```bash
# Set domain for og:* and twitter:* meta tags
xlog -og.domain "example.com"
```

#### ActivityPub Integration

For Fediverse/Mastodon integration:

```bash
xlog -activitypub.username "yourusername" \
     -activitypub.domain "example.com" \
     -activitypub.summary "My knowledge base" \
     -activitypub.icon "/public/avatar.png" \
     -activitypub.image "/public/cover.png"
```

#### Disqus Comments

```bash
xlog -disqus "your-disqus-domain.disqus.com"
```

#### Twitter Card Integration

```bash
xlog -twitter.username "@yourhandle"
```

### Other Flags

#### `-sitemap.domain` - Sitemap Generation

```bash
xlog -sitemap.domain "example.com"
```

#### `-csrf-cookie` - CSRF Cookie Name

```bash
xlog -csrf-cookie "custom_csrf"  # default: "xlog_csrf"
```

#### `-sql-table.threshold` - SQL Query Threshold

Enable SQL queries on tables with more than N rows:

```bash
xlog -sql-table.threshold 100  # default
```

## Common Workflows

### Local Writing Workflow

1. Start XLog server:
   ```bash
   xlog
   ```

2. Open browser to http://localhost:3000

3. Click "Edit" on any page (opens in your configured editor)

4. Edit the markdown file and save

5. Browser automatically reloads with changes

See [Workflow Guide](Workflow.md) for detailed explanation.

### Static Site Deployment

1. Write and preview locally:
   ```bash
   xlog
   ```

2. Generate static site:
   ```bash
   xlog -build ./public
   ```

3. Deploy `./public` directory to your hosting service

### Multi-Site Management

Run multiple XLog instances on different ports:

```bash
# Personal wiki on port 3000
xlog -source ~/wiki -bind :3000 &

# Work notes on port 3001
xlog -source ~/work-notes -bind :3001 &

# Learning journal on port 3002
xlog -source ~/learning -bind :3002 &
```

## Hot-Reload Explanation

XLog watches your markdown files for changes. When you save a file in your editor:

1. XLog detects the file modification
2. Re-renders the markdown to HTML
3. Sends update to browser via WebSocket
4. Browser reloads the page automatically

This happens in milliseconds, creating a seamless writing experience.

## Configuration Files

XLog doesn't use configuration files. All settings are command-line flags. For convenience, create shell aliases or scripts:

```bash
# Add to ~/.bashrc or ~/.zshrc
alias my-wiki='xlog -source ~/wiki -editor "code" -sitename "My Wiki"'
alias work-notes='xlog -source ~/work -editor "vim" -sitename "Work Notes"'
```

## Environment Variables

XLog respects these environment variables:

- `EDITOR` - Default editor (used if `-editor` not specified)
- `PORT` - Server port (overridden by `-bind`)

## Troubleshooting

### Editor Not Opening

Check your editor configuration:

```bash
# Test if editor command works
code test.md  # or vim test.md, emacs test.md, etc.

# Set editor explicitly
xlog -editor "code"
```

### Hot-Reload Not Working

1. Check browser console for WebSocket errors
2. Ensure file is saved (not just open in editor)
3. Try force refresh (Ctrl+Shift+R or Cmd+Shift+R)

### Port Already in Use

Change the port:

```bash
xlog -bind :3001
```

## Next Steps

- [Workflow Guide](Workflow.md) - Understand the editor + preview workflow
- [Extensions](extensions.md) - Learn about XLog's 37 built-in extensions
- [Installation](Installation.md) - Editor setup and configuration

## See All Flags

For the complete list of flags, run:

```bash
xlog -help
```
