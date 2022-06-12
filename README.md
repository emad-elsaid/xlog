XLog
=========

<p align="center"><img width="256" src="assets/logo.png" /></p>

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

## Logo

[Cassette tape icons created by Good Ware - Flaticon](https://www.flaticon.com/free-icons/cassette-tape)

## Screenshots

![](/public/f583bb0cbdf12641666e6f10b26171f61caee3330a68eb825ecbf77eab0227bd.png)

![](/public/e070d8b44a069a6b7336d9be02da4be9020d381aabeb073aa3399df71ca0492b.png)

![](/public/5a182dba45298d4bbc837a6d719c5c194c00f0fce8be363e33f85e9f7f849903.png)

![](/public/7fb69b1749f666f57bcc4044496efd943a84e0f2330cdd724177b4a225baef38.png)

![](/public/bbf6f8be4d374a33338938aeefe70a2f77eb841af8e77b6b79e2659b872e3933.png)

![](/public/9e732a37a60ea4e75c66e5acee8eb493e58d69f52c7121970f4ca24a8f69c8bd.png)

![](/public/8428b0390a330ebc3a815a9efac7b8b6a7b9891f772e6bec5ab73c07d0b9bda4.png)

![](/public/47ef8360209b86693a93794abafb5cccde08499ced4bd5250b77af1e0fce70cd.png)
