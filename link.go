package xlog

// Link represent link for the user interface, default theme renders it in the
// sidebar
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
