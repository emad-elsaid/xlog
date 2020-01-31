XLOG
=========
## What is it?
This about it in the following terms:
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

## Short codes

- Create a directory in your blog called `shortcodes`
- Any executable script in this directory will have be a short code
- Short code syntax is `{scriptName}script input here{/scriptName}`
- The short code part will be replaced by the script output
- The content in the short code will be the first argument of the script
- Scripts can be in any language as long as it's an executable file
- Checkout the `examples/shortcodes` directory for an example of short code scripts

## Installation

### Docker

You can run xlog on port 7000 serving current directory using the following docker command

```
docker run -it --rm -p 7000:7000 -v $PWD:/srv/ emadelsaid/xlog
```

## Usage

```
Usage of xlog:
  -bind string
        IP and port to bind the web server to (default "0.0.0.0:7000")
  -source string
        Directory that will act as a storage (default ".")
```
