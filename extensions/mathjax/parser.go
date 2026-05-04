package mathjax

// MathJax is based on Goldmark-MathJax extension
// https://github.com/litao91/goldmark-mathjax

import (
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

type inlineMathParser struct{}

func (s *inlineMathParser) Trigger() []byte { return []byte{'$'} }

func (s *inlineMathParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, startSegment := block.PeekLine()
	opener := countOpenerDollars(line)
	block.Advance(opener)
	savedLine, savedPos := block.Position()
	node := &InlineMath{}

	if !s.findClosingDelimiter(node, block, startSegment, opener, savedLine, savedPos) {
		return ast.NewTextSegment(startSegment.WithStop(startSegment.Start + opener))
	}

	s.trimHalfSpaces(node, block)
	return node
}

func countOpenerDollars(line []byte) int {
	opener := 0
	for ; opener < len(line) && line[opener] == '$'; opener++ {
	}
	return opener
}

func (s *inlineMathParser) findClosingDelimiter(node *InlineMath, block text.Reader, startSegment text.Segment, opener int, savedLine int, savedPos text.Segment) bool {
	for {
		line, seg := block.PeekLine()
		if line == nil {
			block.SetPosition(savedLine, savedPos)
			return false
		}
		if s.processLine(node, block, line, seg, opener) {
			return true
		}
		if !util.IsBlank(line) {
			node.AppendChild(node, ast.NewRawTextSegment(seg))
		}
		block.AdvanceLine()
	}
}

func (s *inlineMathParser) processLine(node *InlineMath, block text.Reader, line []byte, seg text.Segment, opener int) bool {
	for i := 0; i < len(line); i++ {
		if line[i] == '$' {
			oldi := i
			for ; i < len(line) && line[i] == '$'; i++ {
			}
			closure := i - oldi
			if closure == opener && (i+1 >= len(line) || line[i+1] != '$') {
				closingSeg := seg.WithStop(seg.Start + i - closure)
				if !closingSeg.IsEmpty() {
					node.AppendChild(node, ast.NewRawTextSegment(closingSeg))
				}
				block.Advance(i)
				return true
			}
		}
	}
	return false
}

func (s *inlineMathParser) trimHalfSpaces(node *InlineMath, block text.Reader) {
	if node.IsBlank(block.Source()) {
		return
	}

	firstSeg := node.FirstChild().(*ast.Text).Segment
	lastSeg := node.LastChild().(*ast.Text).Segment

	shouldTrim := !firstSeg.IsEmpty() && block.Source()[firstSeg.Start] == ' ' &&
		!lastSeg.IsEmpty() && block.Source()[lastSeg.Stop-1] == ' '

	if shouldTrim {
		t := node.FirstChild().(*ast.Text)
		t.Segment = t.Segment.WithStart(t.Segment.Start + 1)
		t = node.LastChild().(*ast.Text)
		t.Segment = t.Segment.WithStop(t.Segment.Stop - 1)
	}
}

type mathJaxBlockParser struct{}

type mathBlockData struct {
	indent int
}

var mathBlockInfoKey = parser.NewContextKey()

func (b *mathJaxBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if line == nil {
		return nil, parser.NoChildren
	}

	pos := pc.BlockOffset()
	if pos == -1 {
		return nil, parser.NoChildren
	}
	if line[pos] != '$' {
		return nil, parser.NoChildren
	}
	i := pos
	for ; i < len(line) && line[i] == '$'; i++ {
	}
	if i-pos < 2 {
		return nil, parser.NoChildren
	}
	pc.Set(mathBlockInfoKey, &mathBlockData{indent: pos})
	node := &MathBlock{}
	return node, parser.NoChildren
}

func (b *mathJaxBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if line == nil {
		return parser.NoChildren
	}

	data := pc.Get(mathBlockInfoKey).(*mathBlockData)
	w, pos := util.IndentWidth(line, 0)
	if w < 4 {
		i := pos
		for ; i < len(line) && line[i] == '$'; i++ {
		}
		length := i - pos
		if length >= 2 && util.IsBlank(line[i:]) {
			reader.Advance(segment.Stop - segment.Start - segment.Padding)
			return parser.Close
		}
	}

	_, padding := util.IndentPositionPadding(line, 0, 0, data.indent)
	nonSpacePos := util.FirstNonSpacePosition(line)
	seg := text.NewSegmentPadding(segment.Start+nonSpacePos, segment.Stop, padding)
	node.Lines().Append(seg)
	reader.AdvanceAndSetPadding(segment.Stop-segment.Start-nonSpacePos-1, padding)
	return parser.Continue | parser.NoChildren
}

func (b *mathJaxBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	pc.Set(mathBlockInfoKey, nil)
}

func (b *mathJaxBlockParser) CanInterruptParagraph() bool { return true }
func (b *mathJaxBlockParser) CanAcceptIndentedLine() bool { return false }
func (b *mathJaxBlockParser) Trigger() []byte             { return nil }
