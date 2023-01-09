# Download latest binary

Github has a release for each xlog version tag. it has binaries built for (Windows, Linux, MacOS) for several architectures. you can download the latest version from this page: https://github.com/emad-elsaid/xlog/releases/latest

# Using Go

```bash
go install github.com/emad-elsaid/xlog/cmd/xlog@latest
```

# From source

```bash
git clone git@github.com:emad-elsaid/xlog.git
cd xlog
go run ./cmd/xlog # to run it
go install ./cmd/xlog # to install it to Go bin.
```

# Arch Linux (AUR)

* Xlog is published to AUR: https://aur.archlinux.org/packages/xlog-git
* Using `yay` for example:

```bash
yay -S xlog-git
```

# From source with docker-compose

```bash
git clone git@github.com:emad-elsaid/xlog.git
cd xlog
docker-composer build
docker-composer run
```

```info
Xlog container attach `~/.xlog` as a volume and will write pages to it.
```

# Docker

Releases are packaged as docker images and pushed to GitHub 

```bash
docker pull ghcr.io/emad-elsaid/xlog:latest
docker run -p 3000:3000 -v ~/.xlog:/files ghcr.io/emad-elsaid/xlog:latest
```