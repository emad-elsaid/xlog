# XLog Workflow Guide

This guide explains how XLog works and how to integrate it into your daily workflow.

---

## Understanding the XLog Workflow

XLog follows a **desktop-first editing workflow** that integrates with your existing text editor. Here's what makes it unique:

```
┌─────────────────────────────────────────────────────────┐
│  Your Text Editor (Vim, VS Code, Emacs, etc)          │
│  ↓ Edit & Save                                          │
│  Markdown Files (Filesystem)                            │
│  ↓ XLog watches for changes                             │
│  XLog Server (Live Preview)                             │
│  ↓ Hot-reload                                            │
│  Browser (Instant Preview)                              │
└─────────────────────────────────────────────────────────┘
```

**Key Insight:** XLog doesn't replace your editor—it enhances it with live preview and knowledge-base features.

---

## How Editing Works

### The "Edit" Button

When you click the "Edit" button in your browser:

1. **XLog sends a command** to open the file
2. **Your configured editor launches** with the file open
3. **You edit in your editor** (not in the browser)
4. **You save the file** (Ctrl+S or :w)
5. **XLog detects the change** and reloads the browser

**Important:** There is NO in-browser editing. This is intentional—your editor is more powerful than any web-based solution.

### Configuring Your Editor

XLog uses the `-editor` flag or `$EDITOR` environment variable to know which editor to open.

**Default behavior:**
```bash
xlog  # Uses $EDITOR if set, otherwise tries common editors
```

**Specify editor explicitly:**
```bash
xlog -editor "code"           # VS Code
xlog -editor "vim"            # Vim
xlog -editor "emacs"          # Emacs
xlog -editor "subl -w"        # Sublime Text (wait for close)
xlog -editor "code --wait"    # VS Code (wait for close)
```

**Set permanently in shell:**
```bash
# Add to ~/.bashrc or ~/.zshrc
export EDITOR="code --wait"
```

---

## Editor-Specific Setup

### VS Code

```bash
# Option 1: Set globally
export EDITOR="code --wait"

# Option 2: Pass flag
xlog -editor "code --wait"
```

**Recommended VS Code extensions:**
- Markdown All in One
- Markdown Preview Enhanced
- markdownlint

### Vim

```bash
# Already works if $EDITOR is set
export EDITOR="vim"

# Or use Neovim
export EDITOR="nvim"
```

**Recommended Vim plugins:**
- vim-markdown
- goyo.vim (distraction-free)
- vimwiki (if you want Vim-based wiki too)

### Emacs

```bash
# Use emacsclient for faster startup
export EDITOR="emacsclient -n -a emacs"

# This is XLog's default if no $EDITOR is set
```

**Recommended Emacs packages:**
- markdown-mode
- org-mode (for mixed workflows)

### Sublime Text

```bash
export EDITOR="subl -w"  # -w waits for file close
```

### Other Editors

**Atom:**
```bash
export EDITOR="atom --wait"
```

**Nano:**
```bash
export EDITOR="nano"
```

**Custom GUI editor:**
```bash
# Replace with your editor's command
export EDITOR="/path/to/editor"
```

---

## Daily Workflow Examples

### Workflow 1: Quick Note Taking

```bash
# Terminal 1: Start XLog
cd ~/notes
xlog

# Terminal 2: Create new note
echo "# Meeting Notes - 2026-05-07" > meetings/daily-standup.md

# Browser: Navigate to http://localhost:3000/meetings/daily-standup
# Click "Edit" → Editor opens → Add content → Save
# Browser auto-refreshes with changes
```

### Workflow 2: Research Journal

```bash
# Organize by topic with hashtags
# In your editor, create: research/quantum-computing.md

---
# Quantum Computing Basics

#physics #computing #research

Quantum computers use qubits instead of bits...

Related: [[classical-computing]] [[quantum-entanglement]]
---

# XLog automatically:
# - Creates backlinks from referenced pages
# - Indexes hashtags for browsing
# - Adds to search index
```

### Workflow 3: Documentation Site

```bash
# Project structure
my-docs/
  index.md           # Homepage
  getting-started.md
  api/
    authentication.md
    endpoints.md
  public/
    logo.png
    diagrams/

# Develop with live preview
xlog

# Deploy as static site
xlog -build output/
# Upload output/ to GitHub Pages, Netlify, etc
```

---

## Hot-Reload Explained

XLog watches your filesystem for changes and automatically refreshes your browser. Here's what triggers a reload:

**Triggers hot-reload:**
- Saving a markdown file
- Creating a new file
- Deleting a file
- Renaming a file (via file operations extension)
- Uploading an image

**Does NOT trigger reload:**
- Cursor movement in editor
- Unsaved changes
- Changes to non-markdown files (unless configured)

**How it works:**
1. XLog uses filesystem watchers (inotify on Linux)
2. When a file changes, XLog invalidates caches
3. XLog sends reload signal to browser via WebSocket
4. Browser requests fresh page content
5. You see updated content instantly

**Performance tip:** XLog only watches markdown files by default. If you have thousands of files, this stays fast.

---

## Git Integration

XLog is Git-native, meaning it works perfectly with version control.

### Basic Git Workflow

```bash
# Initialize Git in your notes directory
cd ~/notes
git init

# XLog ignores .git automatically
xlog

# Add and commit as you write
git add .
git commit -m "Add quantum computing notes"

# Push to remote backup
git remote add origin git@github.com:yourname/notes.git
git push -u origin main
```

### Advanced: Auto-Commit on Save

You can create a script that auto-commits after edits:

```bash
#!/bin/bash
# save-and-commit.sh

# Watch for changes and auto-commit
while inotifywait -r -e modify,create,delete ~/notes; do
    cd ~/notes
    git add .
    git commit -m "Auto-commit: $(date)"
done
```

**Note:** Auto-commit is optional. Many users prefer manual commits for meaningful messages.

### Branch-Based Drafts

```bash
# Create draft branch
git checkout -b drafts

# XLog shows same content, just different Git branch
xlog

# Publish by merging to main
git checkout main
git merge drafts
```

---

## Multi-Device Sync

### Option 1: Git + Remote Repo

```bash
# Device 1: Push changes
cd ~/notes
git add .
git commit -m "Update notes"
git push

# Device 2: Pull changes
cd ~/notes
git pull
# XLog automatically detects changes and reloads
```

### Option 2: Obsidian + XLog

Many users use Obsidian for mobile editing and XLog for publishing:

```bash
# Point XLog to Obsidian vault
xlog -source ~/Obsidian/MyVault

# Edit in Obsidian on phone/tablet
# XLog picks up changes automatically
# Publish static site from XLog
```

### Option 3: Syncthing / Dropbox

```bash
# Sync folder across devices
# XLog watches for filesystem changes
# Works with any sync solution
```

---

## Troubleshooting

### Editor Doesn't Open

**Problem:** Clicking "Edit" does nothing

**Solutions:**
1. Check if editor is in PATH: `which code` or `which vim`
2. Set editor explicitly: `xlog -editor "vim"`
3. Check XLog logs for errors
4. Verify file permissions

### Changes Don't Reload

**Problem:** Saved file but browser doesn't update

**Solutions:**
1. Check browser console for WebSocket errors
2. Refresh browser manually (Ctrl+R)
3. Restart XLog server
4. Check file was actually saved (some editors delay writes)

### Permission Errors

**Problem:** Can't edit files

**Solutions:**
1. Check file ownership: `ls -la`
2. Fix permissions: `chmod 644 *.md`
3. Ensure XLog has write access to directory

---

## Best Practices

### Organization

**Use clear folder structure:**
```
notes/
  index.md              # Homepage
  daily/                # Daily notes
  projects/             # Project-specific notes
  reference/            # Reference material
  public/               # Images, assets
```

**Use frontmatter for metadata:**
```yaml
---
title: My Research Paper
tags: [research, physics]
date: 2026-05-07
---
```

### Linking

**Use descriptive page names:**
- ✅ `[[quantum-computing-basics]]`
- ❌ `[[qc1]]`

**Create hub pages:**
- Topic overview pages that link to related notes
- XLog shows backlinks automatically

### Performance

**For large knowledge bases (1000+ pages):**
- Split into subdirectories by topic
- Use `.gitignore` for files you don't need in Git
- Run `xlog -readonly` if you don't need editing

---

## Summary

**XLog workflow = Your Editor + Live Preview + Git**

- ✅ Edit in your powerful desktop editor
- ✅ See changes instantly in browser
- ✅ Version control with Git
- ✅ No lock-in (just markdown files)
- ✅ Works offline

**Not:**
- ❌ Browser-based editing
- ❌ Cloud storage required
- ❌ Proprietary format
- ❌ Database required

This workflow gives you the best of both worlds: powerful editing tools and modern preview experience.

---

## Next Steps

- [Installation Guide](Installation)
- [Extensions](extensions)
- [Creating a Static Site](/tutorials/Creating a site)
- [GitHub Discussions](https://github.com/emad-elsaid/xlog/discussions)
