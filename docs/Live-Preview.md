# Live Preview & Hot-Reload

XLog's live preview with hot-reload creates a seamless writing experience: edit in your text editor, save, and instantly see changes in your browser. No manual refresh needed.

## How It Works

### The Workflow

1. **Start XLog server**
   ```bash
   xlog
   # => Server running at http://localhost:3000
   ```

2. **Open browser** to http://localhost:3000

3. **Click "Edit"** on any page
   - Opens the markdown file in your configured text editor
   - Editor and browser are now synced

4. **Edit and save** in your text editor
   - XLog detects the file change
   - Automatically re-renders the markdown
   - Browser refreshes with updated content

5. **Repeat** - Continuous edit-save-preview loop

No build step, no manual refresh, no context switching.

## Technical Details

### File Watching

XLog monitors markdown files for changes:

```
1. You save file in editor
2. OS notifies XLog of file modification
3. XLog re-renders markdown → HTML
4. WebSocket sends update to browser
5. Browser reloads page
```

This happens in **milliseconds** - feels instant.

### WebSocket Connection

Browser maintains persistent connection to XLog server:

```
Browser ←─ WebSocket ─→ XLog Server
              │
              ├─ File change detected
              └─ Sends reload signal
```

When connection drops (server restarts), browser auto-reconnects.

### Smart Reloading

XLog reloads intelligently:
- **Only affected pages** reload, not entire site
- **Scroll position** attempts to preserve
- **Form state** maintained when possible

## Setting Up Your Editor

### VS Code

**Open from XLog:**
```bash
xlog -editor "code"
```

**Or set default:**
```bash
export EDITOR="code --wait"
```

**VS Code benefits:**
- Side-by-side: editor on left, browser on right
- Markdown preview disabled (use XLog's live preview instead)
- Git integration for version control

### Vim

```bash
export EDITOR="vim"
xlog
```

**Vim workflow:**
- Terminal on left (running XLog)
- Browser on right
- Edit in terminal, see results in browser
- `:w` to save triggers instant preview

### Emacs

```bash
export EDITOR="emacs"
# or for daemon:
export EDITOR="emacsclient -n -a emacs"
xlog
```

**Emacs benefits:**
- Org-mode support (with Pandoc extension)
- Powerful keyboard navigation
- Built-in Git (Magit)

### Sublime Text

```bash
export EDITOR="subl -w"
xlog
```

### Neovim

```bash
export EDITOR="nvim"
xlog
```

**Neovim advantages:**
- Fast startup
- Lua configuration
- Terminal multiplexer (tmux) integration

## Optimal Screen Layout

### Dual Monitor Setup

```
┌─────────────────┐  ┌─────────────────┐
│                 │  │                 │
│   Text Editor   │  │    Browser      │
│                 │  │  (XLog preview) │
│   (markdown)    │  │                 │
│                 │  │                 │
└─────────────────┘  └─────────────────┘
Monitor 1              Monitor 2
```

Edit on one screen, preview on another.

### Single Monitor (Split)

```
┌──────────────────────────────┐
│ Editor         │  Browser    │
│                │             │
│ # Title        │  Title      │
│                │             │
│ Content...     │  Rendered   │
│                │  content    │
└──────────────────────────────┘
```

50/50 split or 60/40 editor/browser.

### Laptop Workflow

Use virtual desktops/workspaces:

```
Desktop 1: Full-screen editor
Desktop 2: Full-screen browser
```

Swipe between them (macOS/Linux gesture, Windows virtual desktops).

## Writing Workflow

### Basic Flow

```
1. Open page in browser
2. Click "Edit" (opens editor)
3. Write content
4. Save (Ctrl+S / Cmd+S)
5. Instantly see rendered result
6. Continue writing
```

### Rapid Iteration

Live preview enables fast iteration:

```
Type → Save → Preview → Adjust → Save → Preview
```

No mental context switch - stay in writing flow.

### Markdown Preview

XLog renders exactly what your readers will see:
- Links
- Images  
- Code blocks
- Tables
- Custom extensions

What you see is what readers get.

## Use Cases

### 1. Blog Post Writing

```
1. Create new post: `echo "# My Post" > blog/new-post.md`
2. Open in browser, click Edit
3. Write content with immediate feedback
4. See formatting, links, images in real-time
5. Publish when satisfied
```

### 2. Documentation

```
1. Edit API docs
2. See code examples rendered with syntax highlighting
3. Verify links work
4. Check table formatting
5. Save knowing it's correct
```

### 3. Learning Notes

```
1. Write notes during course/video
2. Save frequently to see structure
3. Add links to other notes  
4. Verify everything renders correctly
```

### 4. Digital Gardening

```
1. Rough draft ideas quickly
2. See how they look rendered
3. Add backlinks to related pages
4. Refine over time with instant feedback
```

## Performance Tips

### 1. Keep Browser Tab Open

Don't close the preview tab - leave it open and switch between applications. This maintains WebSocket connection for fastest reload.

### 2. Use Fast Editor

Slow editors delay the write:

```
Vim/Neovim      ← Fast save
VS Code         ← Fast save
Emacs           ← Fast save
Heavy IDEs      ← May lag
```

### 3. Disable Editor Preview

Disable your editor's markdown preview - use XLog's instead:

**VS Code:**
```json
{
  "markdown.preview.breaks": false,
  "markdown.preview.doubleClickToSwitchToEditor": false
}
```

Less resource usage, faster workflow.

### 4. Large Files

For very large pages (5000+ lines):
- Consider splitting into smaller pages
- Or tolerate slight delay (still <1 second)

## Troubleshooting

### Changes Not Appearing

**1. Check file was saved**
- Look for saved indicator in editor
- Try saving again (Ctrl+S / Cmd+S)

**2. Check WebSocket connection**
- Open browser console (F12)
- Look for WebSocket errors
- Refresh page to reconnect

**3. Clear browser cache**
- Hard refresh: Ctrl+Shift+R (Windows/Linux)
- Or: Cmd+Shift+R (macOS)

**4. Restart XLog server**
```bash
Ctrl+C  # Stop server
xlog    # Restart
```

### Page Not Updating After Edit

**1. Wrong file open?**
- Verify you're editing the correct markdown file
- Check file path matches page URL

**2. Syntax error in markdown?**
- Check XLog terminal for error messages
- Fix any markdown syntax issues

**3. File permissions?**
```bash
ls -la page.md  # Check if file is writable
```

### Slow Reloads

**1. Large knowledge base?**
- XLog may rebuild index (first load only)
- Subsequent saves should be fast

**2. Slow disk I/O?**
- Check disk performance
- Consider SSD for better experience

**3. Many browser extensions?**
- Disable extensions temporarily
- They can slow page reloads

## Hot-Reload vs. Static Build

### During Writing (Hot-Reload)

```bash
xlog
```

- Instant preview
- No build step
- Perfect for writing

### For Deployment (Static Build)

```bash
xlog -build ./public
```

- Generates static HTML
- Deploy to hosting
- No live server needed

Use hot-reload while writing, static build when deploying.

## Advanced Workflows

### Multiple Pages

Edit multiple pages simultaneously:

```
1. Open Page A in browser
2. Click Edit (opens in editor tab 1)
3. Open Page B in new browser tab
4. Click Edit (opens in editor tab 2)
5. Switch between editor tabs
6. Each save reloads corresponding browser tab
```

### Git Integration

Combine live preview with version control:

```
1. Edit page with live preview
2. Save when satisfied
3. Git: add, commit changes
4. Continue editing

Terminal: git add page.md && git commit -m "Add section on X"
```

### Collaborative Workflow

While XLog doesn't support real-time collaboration, you can:

```
1. Person A writes, commits, pushes
2. Person B pulls, reviews in live preview
3. Person B edits, commits, pushes
4. Repeat

Use: Git for coordination, live preview for review
```

## Comparison to Other Tools

### XLog vs. Static Site Generators

**Hugo, Jekyll, 11ty:**
```
Edit → Save → Run build command → Refresh browser
```
Manual, slow iteration.

**XLog:**
```
Edit → Save → (automatic)
```
Zero-friction workflow.

### XLog vs. Obsidian

**Obsidian:**
- ✅ Built-in live preview
- ✅ No separate browser needed
- ❌ Desktop-only (can't share preview with others)

**XLog:**
- ✅ Live preview in browser
- ✅ Shareable preview URL (localhost)
- ✅ Web-native (works anywhere)
- ✅ Keeps editor choice separate

Both have excellent live preview - XLog's is web-based.

### XLog vs. Notion

**Notion:**
- ✅ Real-time preview (WYSIWYG)
- ❌ Proprietary editor (can't use Vim/Emacs)
- ❌ Cloud-based (requires internet)

**XLog:**
- ✅ Use any editor
- ✅ Fully offline
- ✅ Markdown source (portable)
- ❌ Not WYSIWYG (separate edit/preview)

XLog trades WYSIWYG for editor freedom and offline capability.

## Why Live Preview Matters

### 1. Immediate Feedback

See results instantly:
- Does this link work?
- Is this image correct?
- Does the formatting look right?

No guessing - just save and see.

### 2. Flow State

Stay in writing flow:
```
Think → Write → See → Adjust → Write...
```

No context switch to build/refresh breaks concentration.

### 3. Confidence

Know what you're publishing:
```
Preview matches deployed version exactly
```

No surprises after deployment.

### 4. Learning

Understand markdown faster:
```
Try syntax → Save → See result → Learn
```

Perfect for markdown beginners.

## Best Practices

### 1. Save Frequently

```
Write paragraph → Save → Review → Continue
```

Frequent saves = continuous feedback.

### 2. Use Keyboard Shortcuts

```
Ctrl+S / Cmd+S  - Save (triggers reload)
Ctrl+K          - Search
Alt+Tab         - Switch editor ↔ browser
```

Minimize mouse usage for speed.

### 3. Keep Simple Layout

```
Editor + Browser = Perfect

Avoid:
- Multiple windows
- Complex tiling
- Distractions
```

Simple setup = better focus.

### 4. Disable Distractions

While writing:
- Close other browser tabs
- Silence notifications
- Fullscreen mode (editor + preview only)

## Configuration

Hot-reload is provided by the hotreload extension (enabled by default).

To disable (not recommended):
```bash
xlog -disabled-extensions "hotreload"
```

Without hot-reload, you'll need manual browser refresh.

## See Also

- [Workflow](Workflow.md) - Complete editor integration guide
- [Usage](Usage.md) - Command-line options
- [Installation](Installation.md) - Editor setup
- [Why XLog](Why-XLog.md) - Desktop editor philosophy
