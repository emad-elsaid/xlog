XLog
=========

<p align="center"><img src="assets/logo.svg" /></p>

## What is it?
Think about it in the following terms:
- A personal knowledgebase or Digital Garden
- The file system is the storage
- Uses Markdown format
- Has a web interface
- Minimal interface
- Avoids Javascript like the plague

## Features
- Run it in any directory to create a blog
- Go to any path to create a note
- Checkout [index.md](index) for additional features
- Create a `template` note. xlog will use its content to prefix the editor when creating a new note
- Go files starting with `extension-` are self contained removing a file removes the feature (`views/extension/` sometimes will have a view files for the feature)

## Installation

```
go install github.com/emad-elsaid/xlog@latest
```

## Usage

```
Usage of xlog:
  -bind string
        IP and port to bind the web server to (default "127.0.0.1:7000")
  -source string
        Directory that will act as a storage (default ".")
```

Now you can access notes with `localhost:7000`

If you want you can access it with a name like `notes:7000` by adding this line to `/etc/hosts`

```
127.0.0.1       notes
```

Also you can use `notes` instead of `notes:7000` by redirecting traffic from port 80 to port 7000

```
sudo iptables -t nat -I OUTPUT -p tcp -d 127.0.0.1 --dport 80 -j REDIRECT --to-ports 7000
```

so that means you can create a new note in your browser by visiting `notes/note title here`
