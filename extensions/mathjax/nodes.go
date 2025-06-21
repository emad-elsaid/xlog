package mathjax

import (
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/util"
)

type InlineMath struct {
	ast.BaseInline
}

func (n *InlineMath) IsBlank(source []byte) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		text := c.(*ast.Text).Segment
		if !util.IsBlank(text.Value(source)) {
			return false
		}
	}
	return true
}

func (n *InlineMath) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindInlineMath = ast.NewNodeKind("InlineMath")

func (n *InlineMath) Kind() ast.NodeKind { return KindInlineMath }

type MathBlock struct {
	ast.BaseBlock
}

var KindMathBlock = ast.NewNodeKind("MathBLock")

func (n *MathBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

func (n *MathBlock) Kind() ast.NodeKind { return KindMathBlock }
func (n *MathBlock) IsRaw() bool        { return true }
