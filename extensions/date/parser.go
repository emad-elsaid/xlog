package date

import (
	"time"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&dateParser{}, 999),
	))
}

type dateParser struct{}

func (s *dateParser) Trigger() []byte {
	return []byte{' '}
}

var (
	datePatterns = []string{
		`2006-1-2`,

		`2006-January-2`,
		`2006/January/2`,
		`2006\January\2`,

		`2006-Jan-2`,
		`2006/Jan/2`,
		`2006\Jan\2`,

		`2-January-2006`,
		`2/January/2006`,
		`2\January\2006`,

		`2-Jan-2006`,
		`2/Jan/2006`,
		`2\Jan\2006`,

		`Jan-2-2006`,
		`Jan/2/2006`,
		`Jan\2\2006`,

		`January-2-2006`,
		`January/2/2006`,
		`January\2\2006`,
	}
)

func (s *dateParser) Parse(parent ast.Node, reader text.Reader, pc parser.Context) ast.Node {
	l, _ := reader.PeekLine()
	if len(l) < 2 {
		return nil
	}

	advance := 0

	if l[0] == ' ' {
		advance++
		l = l[1:]
	}

	space := len(l)
	for i, b := range l {
		if b == ' ' {
			space = i
			break
		}
	}

	advance += space
	l = l[:space]

	for _, pattern := range datePatterns {
		t, err := time.Parse(pattern, string(l))
		if err == nil {
			reader.Advance(advance)
			return &DateNode{
				time: t,
			}
		}
	}

	return nil
}
