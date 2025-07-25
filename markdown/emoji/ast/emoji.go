// Package ast defines AST nodes that represetns emoji extension's elements.
package ast

import (
	"fmt"

	gast "github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/emoji/definition"
)

// A Emoji struct represents a emoji of GitHub emoji.
type Emoji struct {
	gast.BaseInline

	ShortName []byte
	Value     *definition.Emoji
}

// Dump implements Node.Dump.
func (n *Emoji) Dump(source []byte, level int) {
	m := map[string]string{
		"ShortName": string(n.ShortName),
		"Value":     fmt.Sprintf("%#v", n.Value),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindEmoji is a NodeKind of the emoji node.
var KindEmoji = gast.NewNodeKind("Emoji")

// Kind implements Node.Kind.
func (n *Emoji) Kind() gast.NodeKind {
	return KindEmoji
}

// NewEmoji returns a new Emoji node.
func NewEmoji(shortName []byte, value *definition.Emoji) *Emoji {
	return &Emoji{
		ShortName: shortName,
		Value:     value,
	}
}
