package images

import "github.com/yuin/goldmark/ast"

var KindColumns = ast.NewNodeKind("ImagesColumns")

type imagesColumns struct {
	ast.Paragraph
}

func (i *imagesColumns) Kind() ast.NodeKind { return KindColumns }
