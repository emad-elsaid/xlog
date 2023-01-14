package autolink_pages

import (
	"strings"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&pageLinkParser{}, 999),
	))
}

type pageLinkParser struct{}

func (_ *pageLinkParser) Trigger() []byte {
	// ' ' indicates any white spaces and a line head
	return []byte{' ', '*', '_', '~', '('}
}

func (s *pageLinkParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	if pc.IsInLinkLabel() {
		return nil
	}

	if autolinkPages == nil {
		UpdatePagesList(nil)
	}

	line, segment := block.PeekLine()
	consumes := 0
	start := segment.Start
	c := line[0]
	// advance if current position is not a line head.
	if c == ' ' || c == '*' || c == '_' || c == '~' || c == '(' {
		consumes++
		start++
		line = line[1:]
	}

	var found Page
	var m int

	for _, p := range autolinkPages {
		if len(line) < len(p.Name()) {
			continue
		}

		// Found a page
		if strings.EqualFold(string(line[0:len(p.Name())]), p.Name()) {
			found = p
			m = len(p.Name())
			break
		}
	}

	if found == nil ||
		(len(line) > m && util.IsAlphaNumeric(line[m])) { // next character is word character
		block.Advance(consumes)
		return nil
	}

	if consumes != 0 {
		s := segment.WithStop(segment.Start + 1)
		ast.MergeOrAppendTextSegment(parent, s)
	}
	consumes += m
	block.Advance(consumes)

	n := ast.NewTextSegment(text.NewSegment(start, start+m))
	link := &PageLink{
		page: found,
	}
	link.AppendChild(link, n)
	return link
}
