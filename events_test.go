package xlog

import (
	"errors"
	"html/template"
	"sync"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

// testPage implements the Page interface for testing events
type testPage struct {
	name    string
	modTime time.Time
}

func (t *testPage) Name() string                   { return t.name }
func (t *testPage) FileName() string               { return t.name + ".md" }
func (t *testPage) Exists() bool                   { return true }
func (t *testPage) Render() template.HTML          { return "" }
func (t *testPage) Content() Markdown              { return "" }
func (t *testPage) Delete() bool                   { return false }
func (t *testPage) Write(Markdown) bool            { return false }
func (t *testPage) ModTime() time.Time             { return t.modTime }
func (t *testPage) AST() ([]byte, ast.Node)        { return nil, nil }

func TestListen(t *testing.T) {
	// Clear any existing handlers to isolate test
	clearPageEvents()

	var called bool
	handler := func(p Page) error {
		called = true
		return nil
	}

	Listen(PageChanged, handler)

	if len(pageEvents[PageChanged]) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(pageEvents[PageChanged]))
	}

	// Trigger to verify handler was registered
	testP := &testPage{name: "test"}
	Trigger(PageChanged, testP)

	if !called {
		t.Error("Handler was not called")
	}
}

func TestListenMultipleHandlers(t *testing.T) {
	clearPageEvents()

	callCount := 0
	handler1 := func(p Page) error {
		callCount++
		return nil
	}
	handler2 := func(p Page) error {
		callCount++
		return nil
	}

	Listen(PageChanged, handler1)
	Listen(PageChanged, handler2)

	if len(pageEvents[PageChanged]) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(pageEvents[PageChanged]))
	}

	testP := &testPage{name: "test"}
	Trigger(PageChanged, testP)

	if callCount != 2 {
		t.Errorf("Expected both handlers to be called, got %d calls", callCount)
	}
}

func TestListenDifferentEvents(t *testing.T) {
	clearPageEvents()

	changedCalled := false
	deletedCalled := false

	Listen(PageChanged, func(p Page) error {
		changedCalled = true
		return nil
	})

	Listen(PageDeleted, func(p Page) error {
		deletedCalled = true
		return nil
	})

	testP := &testPage{name: "test"}

	// Trigger PageChanged
	Trigger(PageChanged, testP)
	if !changedCalled {
		t.Error("PageChanged handler was not called")
	}
	if deletedCalled {
		t.Error("PageDeleted handler should not have been called")
	}

	// Reset flags
	changedCalled = false
	deletedCalled = false

	// Trigger PageDeleted
	Trigger(PageDeleted, testP)
	if changedCalled {
		t.Error("PageChanged handler should not have been called")
	}
	if !deletedCalled {
		t.Error("PageDeleted handler was not called")
	}
}

func TestTriggerNonExistentEvent(t *testing.T) {
	clearPageEvents()

	// Create a non-standard event
	customEvent := PageEvent(999)
	testP := &testPage{name: "test"}

	// Should not panic
	Trigger(customEvent, testP)
}

func TestTriggerWithError(t *testing.T) {
	clearPageEvents()

	expectedErr := errors.New("handler error")
	errorHandler := func(p Page) error {
		return expectedErr
	}

	successCalled := false
	successHandler := func(p Page) error {
		successCalled = true
		return nil
	}

	// Register both handlers
	Listen(PageChanged, errorHandler)
	Listen(PageChanged, successHandler)

	testP := &testPage{name: "test"}
	
	// Trigger should not panic even if a handler returns an error
	Trigger(PageChanged, testP)

	// The second handler should still be called
	if !successCalled {
		t.Error("Second handler should have been called despite first handler error")
	}
}

func TestTriggerPageData(t *testing.T) {
	clearPageEvents()

	var receivedName string
	handler := func(p Page) error {
		receivedName = p.Name()
		return nil
	}

	Listen(PageChanged, handler)

	expectedName := "some/page"
	testP := &testPage{name: expectedName}
	Trigger(PageChanged, testP)

	if receivedName != expectedName {
		t.Errorf("Expected name %s, got %s", expectedName, receivedName)
	}
}

func TestAllPageEvents(t *testing.T) {
	// Verify all defined events can be used
	events := []PageEvent{PageChanged, PageDeleted, PageNotFound}

	clearPageEvents()
	callCounts := make(map[PageEvent]int)

	for _, event := range events {
		evt := event // capture loop variable
		Listen(evt, func(p Page) error {
			callCounts[evt]++
			return nil
		})
	}

	testP := &testPage{name: "test"}

	for _, event := range events {
		Trigger(event, testP)
	}

	for _, event := range events {
		if callCounts[event] != 1 {
			t.Errorf("Event %v: expected 1 call, got %d", event, callCounts[event])
		}
	}
}

func TestListenIdempotency(t *testing.T) {
	clearPageEvents()

	handler := func(p Page) error { return nil }

	// Register same handler multiple times
	Listen(PageChanged, handler)
	Listen(PageChanged, handler)
	Listen(PageChanged, handler)

	// All registrations should be preserved (not deduplicated)
	if len(pageEvents[PageChanged]) != 3 {
		t.Errorf("Expected 3 handlers, got %d", len(pageEvents[PageChanged]))
	}
}

func TestTriggerEmptyHandlerList(t *testing.T) {
	pageEvents = map[PageEvent][]PageEventHandler{
		PageChanged: {},
	}

	testP := &testPage{name: "test"}
	
	// Should not panic with empty handler list
	Trigger(PageChanged, testP)
}

func TestEventHandlerReceivesCorrectPage(t *testing.T) {
	clearPageEvents()

	var (
		receivedPage Page
		mu           sync.Mutex
	)
	handler := func(p Page) error {
		mu.Lock()
		defer mu.Unlock()
		receivedPage = p
		return nil
	}

	Listen(PageDeleted, handler)

	expectedTime := time.Date(2026, 3, 14, 3, 0, 0, 0, time.UTC)
	testP := &testPage{
		name:    "deleted-page",
		modTime: expectedTime,
	}

	Trigger(PageDeleted, testP)

	mu.Lock()
	defer mu.Unlock()
	
	if receivedPage == nil {
		t.Fatal("Handler did not receive page")
	}

	if receivedPage.Name() != "deleted-page" {
		t.Errorf("Expected page name 'deleted-page', got %s", receivedPage.Name())
	}

	if !receivedPage.ModTime().Equal(expectedTime) {
		t.Errorf("Expected mod time %v, got %v", expectedTime, receivedPage.ModTime())
	}
}
