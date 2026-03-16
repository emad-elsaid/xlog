package xlog

import (
	"html/template"
	"testing"
)

// Test priorityList Add and sorting
func TestPriorityListAdd(t *testing.T) {
	pl := &priorityList[string]{}

	// Add items with different priorities
	pl.Add("low", 10.0)
	pl.Add("high", 1.0)
	pl.Add("medium", 5.0)

	// Collect items in order
	var result []string
	for item := range pl.All() {
		result = append(result, item)
	}

	// Should be sorted by priority: high (1.0), medium (5.0), low (10.0)
	expected := []string{"high", "medium", "low"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("at position %d: expected %q, got %q", i, v, result[i])
		}
	}
}

// Test priorityList with equal priorities
func TestPriorityListEqualPriorities(t *testing.T) {
	pl := &priorityList[int]{}

	pl.Add(1, 5.0)
	pl.Add(2, 5.0)
	pl.Add(3, 5.0)

	var result []int
	for item := range pl.All() {
		result = append(result, item)
	}

	// Should maintain stable sort order
	if len(result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result))
	}

	// All items should be present
	found := make(map[int]bool)
	for _, v := range result {
		found[v] = true
	}

	for i := 1; i <= 3; i++ {
		if !found[i] {
			t.Errorf("item %d not found in result", i)
		}
	}
}

// Test priorityList with negative priorities
func TestPriorityListNegativePriorities(t *testing.T) {
	pl := &priorityList[string]{}

	pl.Add("zero", 0.0)
	pl.Add("positive", 10.0)
	pl.Add("negative", -5.0)

	var result []string
	for item := range pl.All() {
		result = append(result, item)
	}

	expected := []string{"negative", "zero", "positive"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("at position %d: expected %q, got %q", i, v, result[i])
		}
	}
}

// Test priorityList All iterator early termination
func TestPriorityListAllEarlyTermination(t *testing.T) {
	pl := &priorityList[int]{}

	for i := 1; i <= 10; i++ {
		pl.Add(i, float32(i))
	}

	var result []int
	for item := range pl.All() {
		result = append(result, item)
		if item == 5 {
			break
		}
	}

	if len(result) != 5 {
		t.Errorf("expected early termination at 5 items, got %d", len(result))
	}

	for i, v := range result {
		if v != i+1 {
			t.Errorf("at position %d: expected %d, got %d", i, i+1, v)
		}
	}
}

// Test RegisterWidget and RenderWidget
func TestRegisterAndRenderWidget(t *testing.T) {
	// Clean up widgets map after test
	defer func() {
		widgets = map[WidgetSpace]*priorityList[WidgetFunc]{}
	}()

	testSpace := WidgetSpace("test_space")
	testPage := NewPage("test")

	// Register widgets with different priorities
	RegisterWidget(testSpace, 10.0, func(p Page) template.HTML {
		return template.HTML("<div>low priority</div>")
	})

	RegisterWidget(testSpace, 1.0, func(p Page) template.HTML {
		return template.HTML("<div>high priority</div>")
	})

	RegisterWidget(testSpace, 5.0, func(p Page) template.HTML {
		return template.HTML("<div>medium priority</div>")
	})

	// Render widgets
	result := RenderWidget(testSpace, testPage)

	expected := template.HTML("<div>high priority</div><div>medium priority</div><div>low priority</div>")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// Test RenderWidget with non-existent space
func TestRenderWidgetNonExistent(t *testing.T) {
	// Clean up widgets map after test
	defer func() {
		widgets = map[WidgetSpace]*priorityList[WidgetFunc]{}
	}()

	testSpace := WidgetSpace("non_existent")
	testPage := NewPage("test")

	result := RenderWidget(testSpace, testPage)

	if result != "" {
		t.Errorf("expected empty HTML for non-existent space, got %q", result)
	}
}

// Test RegisterWidget creates space if not exists
func TestRegisterWidgetCreatesSpace(t *testing.T) {
	// Clean up widgets map after test
	defer func() {
		widgets = map[WidgetSpace]*priorityList[WidgetFunc]{}
	}()

	testSpace := WidgetSpace("new_space")
	testPage := NewPage("test")

	// Register widget to non-existent space
	RegisterWidget(testSpace, 1.0, func(p Page) template.HTML {
		return template.HTML("<div>test</div>")
	})

	// Verify space was created
	if _, ok := widgets[testSpace]; !ok {
		t.Error("expected widget space to be created")
	}

	// Verify widget can be rendered
	result := RenderWidget(testSpace, testPage)
	if result != template.HTML("<div>test</div>") {
		t.Errorf("expected widget to render, got %q", result)
	}
}

// Test WidgetFunc receives correct page
func TestWidgetFuncReceivesPage(t *testing.T) {
	// Clean up widgets map after test
	defer func() {
		widgets = map[WidgetSpace]*priorityList[WidgetFunc]{}
	}()

	testSpace := WidgetSpace("page_test")
	testPage := NewPage("my-test-page")

	var receivedName string
	RegisterWidget(testSpace, 1.0, func(p Page) template.HTML {
		receivedName = p.Name()
		return template.HTML("")
	})

	RenderWidget(testSpace, testPage)

	if receivedName != "my-test-page" {
		t.Errorf("expected page name %q, got %q", "my-test-page", receivedName)
	}
}

// Test multiple widgets in same space accumulate
func TestMultipleWidgetsAccumulate(t *testing.T) {
	// Clean up widgets map after test
	defer func() {
		widgets = map[WidgetSpace]*priorityList[WidgetFunc]{}
	}()

	testSpace := WidgetSpace("accumulate_test")
	testPage := NewPage("test")

	// Register multiple widgets
	for i := 1; i <= 5; i++ {
		priority := float32(i)
		RegisterWidget(testSpace, priority, func(p Page) template.HTML {
			return template.HTML("<span>widget</span>")
		})
	}

	result := RenderWidget(testSpace, testPage)

	// Should have 5 widgets rendered
	expected := template.HTML("<span>widget</span><span>widget</span><span>widget</span><span>widget</span><span>widget</span>")
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// Test predefined widget spaces
func TestPredefinedWidgetSpaces(t *testing.T) {
	tests := []struct {
		space    WidgetSpace
		expected string
	}{
		{WidgetAfterView, "after_view"},
		{WidgetBeforeView, "before_view"},
		{WidgetHead, "head"},
	}

	for _, tt := range tests {
		if string(tt.space) != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, string(tt.space))
		}
	}
}

// Test empty priorityList
func TestEmptyPriorityList(t *testing.T) {
	pl := &priorityList[string]{}

	count := 0
	for range pl.All() {
		count++
	}

	if count != 0 {
		t.Errorf("expected 0 items from empty list, got %d", count)
	}
}

// Test priorityList with single item
func TestPriorityListSingleItem(t *testing.T) {
	pl := &priorityList[string]{}
	pl.Add("only", 1.0)

	var result []string
	for item := range pl.All() {
		result = append(result, item)
	}

	if len(result) != 1 || result[0] != "only" {
		t.Errorf("expected [only], got %v", result)
	}
}

// Test priorityList with float precision
func TestPriorityListFloatPrecision(t *testing.T) {
	pl := &priorityList[string]{}

	pl.Add("a", 1.1)
	pl.Add("b", 1.2)
	pl.Add("c", 1.15)

	var result []string
	for item := range pl.All() {
		result = append(result, item)
	}

	expected := []string{"a", "c", "b"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("at position %d: expected %q, got %q", i, v, result[i])
		}
	}
}
