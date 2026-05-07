# Installation

```callout-info
**Works with Any Text Editor**  
XLog doesn't include a built-in editor. Instead, it works with your favorite text editor (Vim, Emacs, VS Code, Sublime, etc.). You edit markdown files locally, and XLog provides live preview in your browser.
```

## Quick Start (Recommended)

### Using Go

The fastest way to get started:

```bash
go install github.com/emad-elsaid/xlog/cmd/xlog@latest
```

Then create your first note:

```bash
mkdir my-notes
cd my-notes
echo "# Hello World" > index.md
xlog
# => Browse to http://localhost:3000
```

## Download Pre-Built Binary

GitHub releases include binaries for Windows, Linux, and macOS across multiple architectures.

Download the latest version: https://github.com/emad-elsaid/xlog/releases/latest

## Alternative Installation Methods

### From Source

```bash
git clone git@github.com:emad-elsaid/xlog.git
cd xlog
go run ./cmd/xlog # to run it
go install ./cmd/xlog # to install it to Go bin.
```

### Arch Linux (AUR)

Xlog is published to AUR: https://aur.archlinux.org/packages/xlog-git

Using `yay`:

```bash
yay -S xlog-git
```

### Docker

Releases are packaged as docker images and pushed to GitHub Container Registry:

```bash
docker pull ghcr.io/emad-elsaid/xlog:latest
docker run -p 3000:3000 -v ~/.xlog:/files ghcr.io/emad-elsaid/xlog:latest
```

```callout-info
The Docker container mounts `~/.xlog` as a volume and will write pages to it.
```

### Docker Compose (From Source)

```bash
git clone git@github.com:emad-elsaid/xlog.git
cd xlog
docker-compose build
docker-compose run
```

## Editor Setup

After installation, configure your preferred text editor. XLog will use the editor specified in the `-editor` flag or fall back to your `$EDITOR` environment variable.

### VS Code

Set VS Code as your default editor:

```bash
export EDITOR="code"
xlog -editor "code"
```

Or in your shell config (~/.bashrc, ~/.zshrc):

```bash
export EDITOR="code --wait"
```

### Vim

```bash
export EDITOR="vim"
```

### Emacs

```bash
export EDITOR="emacs"
```

### Sublime Text

```bash
export EDITOR="subl -w"
```

### Nano

```bash
export EDITOR="nano"
```

## Verification

After installation, verify XLog is working:

```bash
# Check version
xlog -version

# Start server (creates index.md if missing)
xlog

# Open browser to http://localhost:3000
```

You should see XLog's web interface. Click "Edit" on any page to open it in your configured text editor.

## Troubleshooting

### "Editor not found" Error

If XLog can't find your editor:

1. Set the `EDITOR` environment variable:
   ```bash
   export EDITOR="vim"  # or your preferred editor
   ```

2. Or specify it when running XLog:
   ```bash
   xlog -editor vim
   ```

3. Make sure the editor is in your PATH:
   ```bash
   which vim  # should show the editor path
   ```

### Port 3000 Already in Use

Change the bind address:

```bash
xlog -bind :3001  # Use port 3001 instead
```

### Permission Denied (Linux/macOS)

If installing with Go:

```bash
# Ensure Go bin is in your PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

## Next Steps

- [Workflow Guide](Workflow.md) - Learn the editor + preview workflow
- [Usage](Usage.md) - Command-line flags and options
- [Why XLog](Why-XLog.md) - Understand when to use XLog

## Upgrading

To upgrade to the latest version:

```bash
go install github.com/emad-elsaid/xlog/cmd/xlog@latest
```

Or download the latest binary from the releases page.

See [Upgrading](Upgrading.md) for version-specific migration notes.