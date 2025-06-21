package date

import (
	"fmt"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

var KindDate = ast.NewNodeKind("Date")

type DateNode struct {
	ast.BaseInline
	time time.Time
}

func (d *DateNode) Dump(source []byte, level int) {
	m := map[string]string{
		"value": fmt.Sprintf("%#v", d),
	}
	ast.DumpHelper(d, source, level, m, nil)
}

func (d *DateNode) Kind() ast.NodeKind {
	return KindDate
}
