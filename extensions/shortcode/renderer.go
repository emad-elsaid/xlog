package shortcode

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

type shortCodeRenderer struct{}

func (s *shortCodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindShortCode, s.render)
	reg.Register(KindShortCodeBlock, s.renderBlock)
}

func (s *shortCodeRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	node, ok := n.(*ShortCodeNode)
	if !ok {
		return ast.WalkContinue, nil
	}

	content := source[node.start:node.end]
	output := node.fun.Render(Markdown(content))
	w.Write([]byte(output))

	return ast.WalkContinue, nil
}

func (s *shortCodeRenderer) renderBlock(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	node, ok := n.(*ShortCodeBlock)
	if !ok {
		return ast.WalkContinue, nil
	}

	lines := node.Lines()
	content := ""
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		content += string(line.Value(source))
	}

	output := node.fun.Render(Markdown(content))
	w.Write([]byte(output))

	return ast.WalkContinue, nil
}
