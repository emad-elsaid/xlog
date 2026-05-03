package opengraph

import (
	"flag"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	"github.com/emad-elsaid/xlog"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

var domain string
var twitterUsername string

const descriptionLength = 200

func init() {
	flag.StringVar(&domain, "og.domain", "", "opengraph domain name to be used for meta tags of og:* and twitter:*")
	flag.StringVar(&twitterUsername, "twitter.username", "", "user twitter account @handle. including the @")

	xlog.RegisterExtension(Opengraph{})
}

type Opengraph struct{}

func (Opengraph) Name() string { return "opengraph" }
func (Opengraph) Init() {
	xlog.RegisterWidget(xlog.WidgetHead, 1, opengraphTags)
}

func opengraphTags(p xlog.Page) template.HTML {
	escape := template.JSEscapeString

	title := p.Name()
	if p.Name() == xlog.Config.Index {
		title = xlog.Config.Sitename
	}

	var u url.URL
	u.Scheme = "https"
	u.Host = domain
	u.Path = "/" + title

	URL := u.String()

	var image string
	src, tree := p.AST()
	if imageAST, ok := xlog.FindInAST[*ast.Image](tree); ok && imageAST != nil {
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
		escape(xlog.Config.Sitename),
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

	metaTags := fmt.Sprintf(`
    <meta name="description" content="%s">
`,
		escape(firstParagraph),
	)

	// #nosec G203 -- All variables (title, description, image, etc.) passed through html.EscapeString via escape() function
	return template.HTML(ogTags + twitterTags + metaTags)
}

func rawText(source []byte, n ast.Node, limit int) string {
	if source == nil || n == nil {
		return ""
	}

	out := ""
	_ = ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if n.Kind() == ast.KindText {
			textNode := n.(*ast.Text)
			out += " " + strings.TrimSpace(string(textNode.Segment.Value(source)))
		}

		if len(out) > limit {
			out = out[:limit]
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	return strings.TrimSpace(out)
}
