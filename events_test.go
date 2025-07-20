package xlog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventHandling(t *testing.T) {
	app := newTestApp()
	app.pageEvents = make(map[PageEvent][]PageEventHandler)

	triggered := false
	handler := func(p Page) error {
		triggered = true
		return nil
	}

	app.Listen(PageChanged, handler)
	handlers, exists := app.pageEvents[PageChanged]
	require.True(t, exists, "Expected PageChanged event to be registered")
	require.Len(t, handlers, 1, "Expected 1 handler")

	page := &page{name: "test"}
	app.Trigger(PageChanged, page)
	require.True(t, triggered, "Expected handler to be triggered")
}
