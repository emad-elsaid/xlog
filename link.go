package xlog

func init() {
	RegisterHelper("links", Links)
}

// Link represent link for the user interface, default theme renders it in the
// footer
type Link interface {
	// Icon returns the fontawesome icon class name or emoji
	Icon() string
	// Name returns the link text
	Name() string
	// Link returns the Href property for the link (URL, Path, ...etc)
	Link() string
}

var links = []func(Page) []Link{}

// Register a new links function, should return a list of Links
func RegisterLink(l func(Page) []Link) {
	links = append(links, l)
}

// Links returns a list of links for a Page. it executes all functions
// registered with RegisterLink and collect them in one slice. Can be passed to
// the view to render in the footer or sidebar for example.
func Links(p Page) []Link {
	lnks := []Link{}
	for l := range links {
		lnks = append(lnks, links[l](p)...)
	}
	return lnks
}
