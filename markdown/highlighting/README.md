goldmark-highlighting
=========================

goldmark-highlighting is an extension for the [goldmark](http://github.com/yuin/goldmark) 
that adds syntax-highlighting to the fenced code blocks.

goldmark-highlighting uses [chroma](https://github.com/alecthomas/chroma) as a
syntax highlighter.

Installation
--------------------

```
go get github.com/yuin/goldmark-highlighting/v2
```

Usage
--------------------

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

func main() {
	mdsrc := `
		Title
		=======
		` + "```" + `
		func main() {
		    fmt.Println("ok")
		}
		` + "```" + `
	`

	// Simple usage
	markdown := goldmark.New(
		goldmark.WithExtensions(
			highlighting.Highlighting,
		),
	)
	var buf bytes.Buffer
	if err := markdown.Convert([]byte(mdsrc), &buf); err != nil {
		panic(err)
	}
	title := buf.String()
	fmt.Print(title)

	// Custom configuration
	markdown2 := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
	)
	var buf2 bytes.Buffer
	if err := markdown2.Convert([]byte(mdsrc), &buf2); err != nil {
		panic(err)
	}
	title2 := buf2.String()
	fmt.Print(title2)
}

```

License
--------------------
MIT

Author
--------------------
Yusuke Inuzuka
