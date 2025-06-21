package frontmatter

import (
	"github.com/emad-elsaid/xlog"
)

func init() {
	xlog.RegisterExtension(Frontmatter{})
}

type Frontmatter struct{}

func (Frontmatter) Name() string { return "frontmatter" }
func (Frontmatter) Init() {
	m := New(
		WithStoresInDocument(),
	)

	m.Extend(xlog.MarkdownConverter())
	xlog.RegisterProperty(MetaProperties)
}

type MetaProperty struct {
	NameVal string
	Val     any
}

func (m MetaProperty) Name() string { return m.NameVal }
func (m MetaProperty) Icon() string { return "fa-solid fa-table-list" }
func (m MetaProperty) Value() any   { return m.Val }

func MetaProperties(p xlog.Page) []xlog.Property {
	_, ast := p.AST()
	if ast == nil {
		return nil
	}

	metaData := ast.OwnerDocument().Meta()
	if len(metaData) == 0 {
		return nil
	}

	ps := make([]xlog.Property, 0, len(metaData))
	for k, v := range metaData {
		ps = append(ps, MetaProperty{
			NameVal: k,
			Val:     v,
		})
	}

	return ps
}
