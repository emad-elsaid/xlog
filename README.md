XLog
=========

[![Go Report Card](https://goreportcard.com/badge/github.com/emad-elsaid/xlog)](https://goreportcard.com/report/github.com/emad-elsaid/xlog)
[![GoDoc](https://godoc.org/github.com/emad-elsaid/xlog?status.svg)](https://godoc.org/github.com/emad-elsaid/xlog)



<p align="center"><img width="256" src="public/logo.png" /></p>

Local-first personal knowledge management application with focus on enriching markdown files and surfacing implicit links between pages.

## Core Features

- Uses any directory of markdown files
- Supports Github flavor markdown (GFM)
- One statically compiled binary
- `template.md` content will be used for new pages
- Support for nested directories (although not favored)
- Has a list of tools defined by extensions
- Add image to the start of a page to make it a cover photo
- A web server with a very small footprint thanks to Go.
- Minimal third party dependencies
- The first Emoji used in the page will be considered the icon of the page and displayed beside the title
- Shows task list (Done/Total tasks) beside page link (auto links)

## Extensions

Extensions are defined under `/extensions` sub package. each extension is a subpackage. importing the package in `cmd/xlog/xlog.go` has the side effect of registering the extension hooks. removing the extension from the list of imports will removing the features it provides.

| Extension       | Description                                                                   |
|-----------------|-------------------------------------------------------------------------------|
| Autolink        | Shorten a link string so it wouldn't take unnecessary space                   |
| Autolink pages  | Convert a page name mentions in the middle of text to a link                  |
| Emoji           | Emoji autocomplete while editing                                              |
| File operations | Add a tool item to delete and rename current page                             |
| Hashtags        | Add support for hashtags #hashtag syntax                                      |
| Link preview    | Preview tweets, Facebook posts, youtube videos, Giphy links                   |
| Opengraph       | Adds Opengraph meta tags for title, type, image                               |
| Recent          | Adds an item to sidebar to list all pages ordered by last modified page file. |
| Search          | Full text search                                                              |
| Shortcode       | adds a way for short codes (one line and block)                               |
| Star            | Star pages to pin them to sidebar                                             |
| Upload file     | Add support for upload files, screenshots, audio and camera recording         |
| Versions        | Keeps list of pages older versions                                            |
| Manifest        | adds manifest.json to head tag and output proper JSON value.                  |

## Installation

```
go install github.com/emad-elsaid/xlog/cmd/xlog@latest
```

## Usage

```
Usage of xlog:
  -bind string
        IP and port to bind the web server to (default "127.0.0.1:3000")
  -build string
        Build all pages as static site in this directory
  -index string
        Index file name used as home page (default "index")
  -readonly
        Should xlog hide write operations, read-only means all write operations will be disabled
  -sidebar
        Should render sidebar. (default true)
  -sitename string
        Site name is the name that appears on the header beside the logo and in the title tag (default "XLOG")
  -source string
        Directory that will act as a storage (default current directory)
```

Now you can access xlog with `localhost:3000`

## Generating static site

I used Xlog to generate [my personal blog](https://www.emadelsaid.com/). it uses github workflow to do that [here is an example](https://github.com/emad-elsaid/emad-elsaid.github.io/blob/master/.github/workflows/xlog.yml).

## Overriding Assets

assets is served from `public` directory if it exists in the source directory. otherwise the default assets are served from xlog binary. any file under `public` that has the same name as the ones xlog uses will be used instead of the default files.

## License

Xlog is released under [MIT license](LICENSE)

## Logo

[Cassette tape icons created by Good Ware - Flaticon](https://www.flaticon.com/free-icons/cassette-tape)

## Screenshots

![](/screenshots/285b89e20358e9ea5d1b01893b011665f6282df816983ef1de0d223de698e366.png)![](/screenshots/e9d44ada9ec4190c2ee325df4bbeb789cc67d22dee6bdcdb74393dfa1d8784a3.png)![](/screenshots/75555f02341e1a8ae2775c5f4395b8a52716bd1eeba94cc576c6b6dec5d8c261.png)![](/screenshots/acb69decf484c750f15440c2b39972a03ddaef20509426ed0bb905907fa6154d.png)![](/screenshots/fc52149f89c1e2c1f1b8a352b3eba0743141ed28542a145b1603b3e3f4449db9.png)![](/screenshots/2a8112a513c61a27292753dbbc219eac15f3432b667d38379e79a1d1bb0a629e.png)![](/screenshots/ffa8e45754fca41ff1d76a8e48a296ed13014a2db14eac15eccfea7a83fae1aa.png)
