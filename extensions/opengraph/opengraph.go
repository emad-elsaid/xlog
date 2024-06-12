package opengraph

import (
	"flag"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	. "github.com/emad-elsaid/xlog"

	"github.com/yuin/goldmark/ast"
)

var domain string
var twitterUsername string

const descriptionLength = 200

func init() {
	flag.StringVar(&domain, "og.domain", "", "opengraph domain name to be used for meta tags of og:* and twitter:*")
	flag.StringVar(&twitterUsername, "twitter.username", "", "user twitter account @handle. including the @")
	RegisterWidget(HEAD_WIDGET, 1, opengraphTags)
}

func opengraphTags(p Page) template.HTML {
	escape := template.JSEscapeString

	title := p.Name()
	if p.Name() == INDEX {
		title = SITENAME
	}

	var u url.URL
	u.Scheme = "https"
	u.Host = domain
	u.Path = "/" + title

	URL := u.String()

	var image string
	src, tree := p.AST()
	if imageAST, ok := FindInAST[*ast.Image](tree); ok {
		image = "https://" + domain + string(imageAST.Destination)
	}

	firstParagraph := rawText(src, tree, descriptionLength)

	ogTags := fmt.Sprintf(`
    <meta property="og:site_name" content="%s" />
    <meta property="og:title" content="%s" />
    <meta property="og:description" content="%s" />
    <meta property="og:image" content="%s" />
    <meta property="og:url" content="%s" />
    <meta property="og:type" content="website" />
`,
		escape(SITENAME),
		escape(title),
		escape(firstParagraph),
		escape(image),
		escape(URL),
	)

	twitterTags := fmt.Sprintf(`
    <meta name="twitter:title" content="%s" />
    <meta name="twitter:description" content="%s" />
    <meta name="twitter:image" content="%s" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:creator" content="%s" />
    <meta name="twitter:site" content="%s" />
    <meta name="twitter:image:alt" content="%s" />
`,
		escape(title),
		escape(firstParagraph),
		escape(image),
		escape(twitterUsername),
		escape(twitterUsername),
		escape(title),
	)

	return template.HTML(ogTags + twitterTags)
}

func rawText(source []byte, n ast.Node, limit int) string {
	out := ""
	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() == ast.KindText {
			out += " " + strings.TrimSpace(string(n.(*ast.Text).Text(source)))
		}

		if len(out) > limit {
			out = out[:limit]
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	return strings.TrimSpace(out)
}
