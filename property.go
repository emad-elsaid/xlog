package xlog

// Property represent a piece of information about the current page such as last
// update time, number of versions, number of words, reading time...etc
type Property interface {
	// Icon returns the fontawesome icon class name or emoji
	Icon() string
	// Name returns the link text
	Name() string
}

var props = []func(Page) []Property{defaultProps}

// RegisterProperty registers a function that returns a set of properties for
// the page
func RegisterProperty(a func(Page) []Property) {
	props = append(props, a)
}

// Properties return a list of properties for a page. It executes all functions
// registered with RegisterProperty and collect results in one slice. Can be
// passed to the view to render a page properties
func Properties(p Page) []Property {
	ps := []Property{}
	for pr := range props {
		ps = append(ps, props[pr](p)...)
	}
	return ps
}

type lastUpdateProp struct{ page Page }

func (a lastUpdateProp) Icon() string { return "fa-solid fa-clock" }
func (a lastUpdateProp) Name() string { return ago(a.page.ModTime()) }

func defaultProps(p Page) []Property {
	return []Property{
		lastUpdateProp{p},
	}
}
