package xlog

// Autocompletion defines what character triggeres the autocomplete feature and
// what is the list to display in this case.
type Autocompletion struct {
	StartChar   string
	Suggestions []*Suggestion
}

// Suggestions represent an item in the list of autocomplete menu in the edit page
type Suggestion struct {
	Text        string // The text that gets injected in the editor if this option is choosen
	DisplayText string // The display text for this item in the menu. this can be more cosmetic.
}

// This is a function that returns an auto completer instance. this function
// should be defined by extensions and registered to be executed when rendering
// the edit page
type Autocompleter func() *Autocompletion

// Holds a list of registered autocompleter functions
var autocompletes = []Autocompleter{}

// this function registers an autocompleter function. it should be used by an
// extension to register a new autocompleter function. these functions are going
// to be executed when rendering the edit page.
func Autocomplete(a Autocompleter) {
	autocompletes = append(autocompletes, a)
}
