![](/docs/public/static.png)

#tutorial

Xlog CLI allow for generating static website from source directory. this is how this website is generated.

To generate a static website using Xlog use the `--build` flag with a path as destination for example:

```shell
xlog --build /path/to/output
```

Xlog will build all markdown files to HTML and extract all static files from inside the binary executable file to that destination directory. Then it will terminate.

Building process creates a xlog server instance and request all pages and save it to desk. That allow xlog extensions to define a new handler that renders a page. the page will work in both usecases: local server, static site generation. extensions has to also register the path for build using [`RegisteBuildPage`](https://pkg.go.dev/github.com/emad-elsaid/xlog#RegisterBuildPage) function

While building static site xlog turns on **READONLY** mode.

