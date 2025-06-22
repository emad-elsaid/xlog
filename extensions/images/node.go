package images

import "github.com/emad-elsaid/xlog/markdown/ast"

var KindColumns = ast.NewNodeKind("ImagesColumns")

type imagesColumns struct {
	ast.Paragraph
}

func (i *imagesColumns) Kind() ast.NodeKind { return KindColumns }
