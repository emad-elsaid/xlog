package xlog

type (
	// a type used to define events to be used when the page is manipulated for
	// example modified, renamed, deleted...etc.
	PageEvent int
	// a function that handles a page event. this should be implemented by an
	// extension and then registered. it will get executed when the event is
	// triggered
	PageEventHandler func(Page) error
)

// List of page events. extensions can use these events to register a function
// to be executed when this event is triggered. extensions that require to be
// notified when the page is created or overwritten or deleted should register
// an event handler for the interesting events.
const (
	PageChanged PageEvent = iota
	PageDeleted
	PageNotFound // user requested a page that's not found
)

// Register an event handler to be executed when PageEvent is triggered.
// extensions can use this to register hooks under specific page events.
// extensions that keeps a cached version of the pages list for example needs to
// register handlers to update its cache
func Listen(e PageEvent, h PageEventHandler) {
	app := GetApp()
	app.Listen(e, h)
}

// Trigger event handlers for a specific page event. page methods use this
// function to trigger all registered handlers when the page is edited or
// deleted for example.
func Trigger(e PageEvent, p Page) {
	app := GetApp()
	app.Trigger(e, p)
}
