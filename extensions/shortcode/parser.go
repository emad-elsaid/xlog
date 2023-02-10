package shortcode

import (
	"strings"

	. "github.com/emad-elsaid/xlog"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(&shortCodeParser{}, 0),
	))
}

const trigger = '/'

type shortCodeParser struct{}

func (s *shortCodeParser) Trigger() []byte {
	return []byte{trigger}
}

func (s *shortCodeParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	l, seg := reader.PeekLine()
	line := string(l)
	if len(line) == 0 || line[0] != trigger {
		return nil, parser.NoChildren
	}

	firstSpace := strings.Index(line, " ")
	if firstSpace == -1 {
		firstSpace = len(line) - 1
	}

	firstWord := line[1:firstSpace]
	var processor ShortCodeFunc
	var ok bool
	if processor, ok = shortcodes[firstWord]; !ok {
		return nil, parser.NoChildren
	}

	reader.AdvanceLine()

	end := seg.Stop
	start := seg.Start + len(firstWord) + 2
	if start > end {
		start = end
	}

	return &ShortCodeNode{
		start: start,
		end:   end,
		fun:   processor,
	}, parser.NoChildren
}

func (s *shortCodeParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

func (s *shortCodeParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {}
func (s *shortCodeParser) CanInterruptParagraph() bool                                { return true }
func (s *shortCodeParser) CanAcceptIndentedLine() bool                                { return false }
