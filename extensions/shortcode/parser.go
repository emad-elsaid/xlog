package shortcode

import (
	"strings"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
)

const trigger = '/'

type shortCodeParser struct{}

func (s *shortCodeParser) Trigger() []byte {
	return []byte{trigger}
}

func (s *shortCodeParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	l, seg := reader.PeekLine()
	line := string(l)
	if len(line) == 0 || line[0] != trigger {
		return nil, parser.Close
	}

	endOfShortcode := strings.IndexAny(line, " \n")
	if endOfShortcode == -1 {
		endOfShortcode = len(line)
	}

	firstWord := line[1:endOfShortcode]
	var processor ShortCode
	var ok bool
	if processor, ok = shortcodes[firstWord]; !ok {
		return nil, parser.Close
	}

	reader.AdvanceLine()

	firstSpace := strings.IndexAny(line, " ")
	if firstSpace == -1 {
		return &ShortCodeNode{
			start: seg.Stop,
			end:   seg.Stop,
			fun:   processor,
		}, parser.Close
	}

	return &ShortCodeNode{
		start: seg.Start + endOfShortcode,
		end:   seg.Stop,
		fun:   processor,
	}, parser.Close
}

func (s *shortCodeParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	return parser.Close
}

func (s *shortCodeParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {}
func (s *shortCodeParser) CanInterruptParagraph() bool                                { return true }
func (s *shortCodeParser) CanAcceptIndentedLine() bool                                { return false }
