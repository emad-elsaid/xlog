package parser

import (
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

type blockquoteParser struct {
}

var defaultBlockquoteParser = &blockquoteParser{}

// NewBlockquoteParser returns a new BlockParser that
// parses blockquotes.
func NewBlockquoteParser() BlockParser {
	return defaultBlockquoteParser
}

func (b *blockquoteParser) process(reader text.Reader) bool {
	line, _ := reader.PeekLine()
	w, pos := util.IndentWidth(line, reader.LineOffset())
	if w > 3 || pos >= len(line) || line[pos] != '>' {
		return false
	}
	pos++
	if pos >= len(line) || line[pos] == '\n' {
		reader.Advance(pos)
		return true
	}
	reader.Advance(pos)
	if line[pos] == ' ' || line[pos] == '\t' {
		padding := 0
		if line[pos] == '\t' {
			padding = util.TabWidth(reader.LineOffset()) - 1
		}
		reader.AdvanceAndSetPadding(1, padding)
	}
	return true
}

func (b *blockquoteParser) Trigger() []byte {
	return []byte{'>'}
}

func (b *blockquoteParser) Open(parent ast.Node, reader text.Reader, pc Context) (ast.Node, State) {
	if b.process(reader) {
		return ast.NewBlockquote(), HasChildren
	}
	return nil, NoChildren
}

func (b *blockquoteParser) Continue(node ast.Node, reader text.Reader, pc Context) State {
	if b.process(reader) {
		return Continue | HasChildren
	}
	return Close
}

func (b *blockquoteParser) Close(node ast.Node, reader text.Reader, pc Context) {
	// nothing to do
}

func (b *blockquoteParser) CanInterruptParagraph() bool {
	return true
}

func (b *blockquoteParser) CanAcceptIndentedLine() bool {
	return false
}
