![](/docs/public/custom.png)

#tutorial

Xlog ships with a CLI that includes the core and all official extensions. There are cases where you need custom set of extensions:

* Maybe you need only the core features without any extension
* Maybe there is an extension that you don't need or want or misbehaving
* Maybe you developed a set of extensions and you want to include them in your installations

Here is how you can build your own custom xlog with features you select.

# Creating a Go module

Create a directory for your custom installation and initialize a go module in it.

```shell
mkdir custom_xlog
cd custom_xlog
go mod init github.com/yourusername/custom_xlog
```

# Main file

Then create a file `xlog.go` for example with the following content

```go
package main

import (
	// Core
	"github.com/emad-elsaid/xlog"

	// All official extensions
	_ "github.com/emad-elsaid/xlog/extensions/all"
)

func main() {
	xlog.Start()
}
```

# Selecting extensions

The previous file is what xlog ships in `cmd/xlog/xlog.go` if you missed up at any point feel free to go back to it and copy it from there. 

If you want to select specific extensions you can replace `extensions/all` line with a list of extensions that you want.

All extensions are imported to [`extensions/all/all.go`](https://github.com/emad-elsaid/xlog/blob/master/extensions/all/all.go). feel free to copy any of them as needed.

You can also import any extensions that you developed at this point.

# Running your custom xlog

Now use Go to run your custom installation 

```shell
go get github.com/emad-elsaid/xlog
go run xlog.go
```

