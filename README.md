XLOG
=========
## What is it?
Think about it in the following terms:
- A personal knowledgebase
- The file system is the storage
- Uses Markdown format
- Has a web interface
- Minimal interface
- Avoids Javascript like the plague

## Features
- Run it in any directory to create a blog
- Go to any path to create a post
- Go to any path + `/edit` to edit a post
- Save a post without content to delete it
- Short codes
- Autolinking text to other posts automatically
- Convert youtube links to embded videos

## Short codes

- Create a directory in your blog called `shortcodes`
- Any executable script in this directory will have be a short code
- Short code syntax is `{scriptName}script input here{/scriptName}`
- The short code part will be replaced by the script output
- The content in the short code will be the STDIN of the script process
- Scripts can be in any language as long as it's an executable file
- Checkout the `examples/shortcodes` directory for an example of short code scripts
- If the first line in the output is `text/markdown` the text will be rendered as markdown, if `text/html` it will be printed as is
- If the first line isn't a mimetype it will be assumed `text/html`

## Installation

### Docker

You can run xlog on port 7000 serving current directory using the following docker command

```
docker run -it --rm -p 7000:7000 -v $PWD:/srv/ emadelsaid/xlog
```

### From source

```
go get github.com/emad-elsaid/xlog/...
```

## Usage

```
Usage of xlog:
  -bind string
        IP and port to bind the web server to (default "0.0.0.0:7000")
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
