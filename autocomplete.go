package xlog

// Autocomplete defines what character triggeres the autocomplete feature and
// what is the list to display in this case.
type Autocomplete interface {
	StartChar() string
	Suggestions() []*Suggestion
}

// Suggestions represent an item in the list of autocomplete menu in the edit page
type Suggestion struct {
	Text        string // The text that gets injected in the editor if this option is chosen
	DisplayText string // The display text for this item in the menu. this can be more cosmetic.
}

// Holds a list of registered autocomplete functions
var autocompletes = []Autocomplete{}

// RegisterAutocomplete registers an autocomplete function. it should be used by an
// extension to register a new autocomplete function. these functions are going
// to be executed when rendering the edit page.
func RegisterAutocomplete(a Autocomplete) {
	autocompletes = append(autocompletes, a)
}
