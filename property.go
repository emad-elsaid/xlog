package xlog

// Property represent a piece of information about the current page such as last
// update time, number of versions, number of words, reading time...etc
type Property interface {
	// Icon returns the fontawesome icon class name or emoji
	Icon() string
	// Name returns the name of the property
	Name() string
	// Value returns the value of the property
	Value() any
}

var propsSources = []func(Page) []Property{defaultProps}

// RegisterProperty registers a function that returns a set of properties for
// the page
func RegisterProperty(a func(Page) []Property) {
	propsSources = append(propsSources, a)
}

// Properties return a list of properties for a page. It executes all functions
// registered with RegisterProperty and collect results in one slice. Can be
// passed to the view to render a page properties
func Properties(p Page) map[string]Property {
	ps := map[string]Property{}
	for _, source := range propsSources {
		for _, pr := range source(p) {
			ps[pr.Name()] = pr
		}
	}

	return ps
}

type lastUpdateProp struct{ page Page }

func (a lastUpdateProp) Icon() string { return "fa-solid fa-clock" }
func (a lastUpdateProp) Name() string { return "modified" }
func (a lastUpdateProp) Value() any   { return ago(a.page.ModTime()) }

func defaultProps(p Page) []Property {
	if p.ModTime().IsZero() {
		return nil
	}

	return []Property{
		lastUpdateProp{p},
	}
}
