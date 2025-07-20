package xlog

import "log/slog"

// PageEvent represents different events that can occur with a page
type PageEvent int

// PageEventHandler is a function that handles a page event
type PageEventHandler func(Page) error

// List of page events
const (
	PageChanged PageEvent = iota
	PageDeleted
	PageNotFound // user requested a page that's not found
)

// Listen registers an event handler
func (app *App) Listen(e PageEvent, h PageEventHandler) {

	if _, ok := app.pageEvents[e]; !ok {
		app.pageEvents[e] = []PageEventHandler{}
	}

	app.pageEvents[e] = append(app.pageEvents[e], h)
}

// Trigger triggers an event
func (app *App) Trigger(e PageEvent, p Page) {
	handlers, ok := app.pageEvents[e]

	if !ok {
		return
	}

	for _, h := range handlers {
		if err := h(p); err != nil {
			slog.Error("Failed to execute handler for event", "event", e, "handler", h, "error", err)
		}
	}
}
