package main

import (
	"fmt"
	"html/template"

	"github.com/yuin/goldmark/ast"
)

func init() {
	WIDGET(HEAD_WIDGET, opengraphTags)
}

func opengraphTags(p *Page, r Request) template.HTML {
	tags := `<meta property="og:type" content="article" />`
	name := p.Name
	if p.Name == "index" {
		name = SITENAME
	}

	tags += fmt.Sprintf(`<meta property="og:title" content="%s%s" />`, p.Emoji(), template.JSEscapeString(name))

	if image, ok := ExtractFirstFromAST[*ast.Image](p.AST(), ast.KindImage); ok {
		tags += fmt.Sprintf(`<meta property="og:image" content="%s" />`, template.JSEscapeString(string(image.Destination)))
	}

	return template.HTML(tags)
}