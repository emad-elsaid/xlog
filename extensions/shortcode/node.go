package shortcode

import (
	"fmt"

	"github.com/yuin/goldmark/ast"
)

var KindShortCode = ast.NewNodeKind("ShortCode")

type ShortCodeNode struct {
	ast.BaseBlock
	start int
	end   int
	fun   ShortCodeFunc
}

func (s *ShortCodeNode) Dump(source []byte, level int) {
	m := map[string]string{
		"value": fmt.Sprintf("%#v", s),
	}
	ast.DumpHelper(s, source, level, m, nil)
}

func (h *ShortCodeNode) Kind() ast.NodeKind {
	return KindShortCode
}
