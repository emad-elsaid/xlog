package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

const (
	testComputedValue = "computed"
)

func TestReference(t *testing.T) {
	tests := []struct {
		name        string
		label       []byte
		destination []byte
		title       []byte
	}{
		{
			name:        "basic reference",
			label:       []byte("link"),
			destination: []byte("https://example.com"),
			title:       []byte("Example Site"),
		},
		{
			name:        "empty title",
			label:       []byte("link"),
			destination: []byte("https://example.com"),
			title:       []byte(""),
		},
		{
			name:        "empty label",
			label:       []byte(""),
			destination: []byte("https://example.com"),
			title:       []byte("Example"),
		},
		{
			name:        "all empty",
			label:       []byte(""),
			destination: []byte(""),
			title:       []byte(""),
		},
		{
			name:        "unicode in label",
			label:       []byte("日本語"),
			destination: []byte("https://example.jp"),
			title:       []byte("Japanese Site"),
		},
		{
			name:        "special characters",
			label:       []byte("my-link_123"),
			destination: []byte("/path/to/file?query=1&foo=bar"),
			title:       []byte("Title with \"quotes\""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ref := NewReference(tc.label, tc.destination, tc.title)

			if string(ref.Label()) != string(tc.label) {
				t.Errorf("Label() = %s, want %s", ref.Label(), tc.label)
			}
			if string(ref.Destination()) != string(tc.destination) {
				t.Errorf("Destination() = %s, want %s", ref.Destination(), tc.destination)
			}
			if string(ref.Title()) != string(tc.title) {
				t.Errorf("Title() = %s, want %s", ref.Title(), tc.title)
			}

			// Test String() method
			str := ref.String()
			if str == "" {
				t.Error("String() returned empty string")
			}
		})
	}
}

func TestIDs_Generate(t *testing.T) {
	tests := []struct {
		name     string
		value    []byte
		kind     ast.NodeKind
		expected string
	}{
		{
			name:     "simple alphanumeric",
			value:    []byte("Hello World"),
			kind:     ast.KindHeading,
			expected: "hello-world",
		},
		{
			name:     "with special characters",
			value:    []byte("Hello, World!"),
			kind:     ast.KindHeading,
			expected: "hello-world",
		},
		{
			name:     "with hyphens and underscores",
			value:    []byte("hello-world_test"),
			kind:     ast.KindHeading,
			expected: "hello-world-test",
		},
		{
			name:     "uppercase conversion",
			value:    []byte("HELLO WORLD"),
			kind:     ast.KindHeading,
			expected: "hello-world",
		},
		{
			name:     "empty heading",
			value:    []byte(""),
			kind:     ast.KindHeading,
			expected: "heading",
		},
		{
			name:     "empty non-heading",
			value:    []byte(""),
			kind:     ast.KindDocument,
			expected: "id",
		},
		{
			name:     "unicode characters",
			value:    []byte("Hello 日本語 World"),
			kind:     ast.KindHeading,
			expected: "hello--world",
		},
		{
			name:     "only special characters",
			value:    []byte("!!!"),
			kind:     ast.KindHeading,
			expected: "heading",
		},
		{
			name:     "leading and trailing spaces",
			value:    []byte("  hello world  "),
			kind:     ast.KindHeading,
			expected: "hello-world",
		},
		{
			name:     "multiple spaces",
			value:    []byte("hello    world"),
			kind:     ast.KindHeading,
			expected: "hello----world",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ids := newIDs()
			result := ids.Generate(tc.value, tc.kind)
			if string(result) != tc.expected {
				t.Errorf("Generate(%q) = %q, want %q", tc.value, result, tc.expected)
			}
		})
	}
}

func TestIDs_GenerateDuplicates(t *testing.T) {
	ids := newIDs()

	// First generation
	id1 := ids.Generate([]byte("test"), ast.KindHeading)
	if string(id1) != "test" {
		t.Errorf("First id = %q, want 'test'", id1)
	}

	// Second generation should append -1
	id2 := ids.Generate([]byte("test"), ast.KindHeading)
	if string(id2) != "test-1" {
		t.Errorf("Second id = %q, want 'test-1'", id2)
	}

	// Third generation should append -2
	id3 := ids.Generate([]byte("test"), ast.KindHeading)
	if string(id3) != "test-2" {
		t.Errorf("Third id = %q, want 'test-2'", id3)
	}
}

func TestIDs_Put(t *testing.T) {
	ids := newIDs()

	// Put a value manually
	ids.Put([]byte("reserved"))

	// Generating the same value should create a numbered variant
	result := ids.Generate([]byte("reserved"), ast.KindHeading)
	if string(result) != "reserved-1" {
		t.Errorf("Generate after Put = %q, want 'reserved-1'", result)
	}
}

func TestContext_GetSet(t *testing.T) {
	key := NewContextKey()
	ctx := NewContext()

	// Initially should be nil
	if val := ctx.Get(key); val != nil {
		t.Errorf("Get() before Set = %v, want nil", val)
	}

	// Set a value
	testValue := "test value"
	ctx.Set(key, testValue)

	// Get should return the value
	if val := ctx.Get(key); val != testValue {
		t.Errorf("Get() after Set = %v, want %v", val, testValue)
	}
}

func TestContext_ComputeIfAbsent(t *testing.T) {
	key := NewContextKey()
	ctx := NewContext()

	callCount := 0
	computeFunc := func() any {
		callCount++
		return testComputedValue
	}

	// First call should compute
	val1 := ctx.ComputeIfAbsent(key, computeFunc)
	if val1 != testComputedValue {
		t.Errorf("First ComputeIfAbsent = %v, want 'computed'", val1)
	}
	if callCount != 1 {
		t.Errorf("Compute function called %d times, want 1", callCount)
	}

	// Second call should not compute
	val2 := ctx.ComputeIfAbsent(key, computeFunc)
	if val2 != testComputedValue {
		t.Errorf("Second ComputeIfAbsent = %v, want 'computed'", val2)
	}
	if callCount != 1 {
		t.Errorf("Compute function called %d times, want 1", callCount)
	}
}

func TestContext_References(t *testing.T) {
	ctx := NewContext()

	// Add references
	ref1 := NewReference([]byte("link1"), []byte("http://example.com"), []byte("Example"))
	ref2 := NewReference([]byte("link2"), []byte("http://test.com"), []byte("Test"))

	ctx.AddReference(ref1)
	ctx.AddReference(ref2)

	// Test Reference lookup
	foundRef, ok := ctx.Reference("link1")
	if !ok {
		t.Error("Reference('link1') not found")
	}
	if foundRef != nil && string(foundRef.Destination()) != "http://example.com" {
		t.Errorf("Reference destination = %s, want 'http://example.com'", foundRef.Destination())
	}

	// Test non-existent reference
	_, ok = ctx.Reference("nonexistent")
	if ok {
		t.Error("Reference('nonexistent') should not be found")
	}

	// Test References list
	refs := ctx.References()
	if len(refs) != 2 {
		t.Errorf("References() length = %d, want 2", len(refs))
	}
}

func TestContext_AddReferenceDuplicate(t *testing.T) {
	ctx := NewContext()

	ref1 := NewReference([]byte("link"), []byte("http://example.com"), []byte("Example"))
	ref2 := NewReference([]byte("link"), []byte("http://different.com"), []byte("Different"))

	ctx.AddReference(ref1)
	ctx.AddReference(ref2)

	// First reference should win
	foundRef, ok := ctx.Reference("link")
	if !ok {
		t.Error("Reference('link') not found")
	}
	if foundRef != nil && string(foundRef.Destination()) != "http://example.com" {
		t.Errorf("Reference destination = %s, want 'http://example.com'", foundRef.Destination())
	}
}

func TestContext_BlockOffsetAndIndent(t *testing.T) {
	ctx := NewContext()

	// Initial values should be -1
	if offset := ctx.BlockOffset(); offset != -1 {
		t.Errorf("Initial BlockOffset = %d, want -1", offset)
	}
	if indent := ctx.BlockIndent(); indent != -1 {
		t.Errorf("Initial BlockIndent = %d, want -1", indent)
	}

	// Set values
	ctx.SetBlockOffset(5)
	ctx.SetBlockIndent(3)

	if offset := ctx.BlockOffset(); offset != 5 {
		t.Errorf("BlockOffset after Set = %d, want 5", offset)
	}
	if indent := ctx.BlockIndent(); indent != 3 {
		t.Errorf("BlockIndent after Set = %d, want 3", indent)
	}
}

func TestContext_Delimiters(t *testing.T) {
	ctx := NewContext()

	// Initially should be nil
	if first := ctx.FirstDelimiter(); first != nil {
		t.Error("FirstDelimiter should be nil initially")
	}
	if last := ctx.LastDelimiter(); last != nil {
		t.Error("LastDelimiter should be nil initially")
	}

	// Create and push a delimiter
	delim := NewDelimiter(true, false, 1, '*', nil)
	ctx.PushDelimiter(delim)

	if first := ctx.FirstDelimiter(); first != delim {
		t.Error("FirstDelimiter should return the pushed delimiter")
	}
	if last := ctx.LastDelimiter(); last != delim {
		t.Error("LastDelimiter should return the pushed delimiter")
	}
}

func TestContext_OpenedBlocks(t *testing.T) {
	ctx := NewContext()

	// Initially should be empty
	blocks := ctx.OpenedBlocks()
	if len(blocks) != 0 {
		t.Errorf("Initial OpenedBlocks length = %d, want 0", len(blocks))
	}

	// LastOpenedBlock on empty should return zero value
	lastBlock := ctx.LastOpenedBlock()
	if lastBlock.Node != nil {
		t.Error("LastOpenedBlock should have nil Node initially")
	}

	// Set blocks
	testBlocks := []Block{
		{Node: ast.NewParagraph(), Parser: nil},
		{Node: ast.NewHeading(1), Parser: nil},
	}
	ctx.SetOpenedBlocks(testBlocks)

	blocks = ctx.OpenedBlocks()
	if len(blocks) != 2 {
		t.Errorf("OpenedBlocks length = %d, want 2", len(blocks))
	}

	lastBlock = ctx.LastOpenedBlock()
	if lastBlock.Node == nil || lastBlock.Node.Kind() != ast.KindHeading {
		t.Error("LastOpenedBlock should return the last heading")
	}
}

func TestContext_IDs(t *testing.T) {
	ctx := NewContext()

	ids := ctx.IDs()
	if ids == nil {
		t.Error("IDs() should not return nil")
	}

	// Test that it works
	id := ids.Generate([]byte("test"), ast.KindHeading)
	if string(id) != "test" {
		t.Errorf("IDs().Generate() = %q, want 'test'", id)
	}
}

func TestContext_WithIDs(t *testing.T) {
	customIDs := newIDs()
	customIDs.Put([]byte("reserved"))

	ctx := NewContext(WithIDs(customIDs))

	// The custom IDs should be used
	result := ctx.IDs().Generate([]byte("reserved"), ast.KindHeading)
	if string(result) != "reserved-1" {
		t.Errorf("Generate with custom IDs = %q, want 'reserved-1'", result)
	}
}

func TestContext_String(t *testing.T) {
	ctx := NewContext()
	str := ctx.String()
	if str == "" {
		t.Error("String() should not return empty string")
	}
}

func TestState_Constants(t *testing.T) {
	// Test that state constants are distinct
	states := []State{None, Continue, Close, HasChildren, NoChildren, RequireParagraph}

	for i, s1 := range states {
		for j, s2 := range states {
			if i != j && s1 == s2 {
				t.Errorf("States at index %d and %d are equal: %v", i, j, s1)
			}
		}
	}
}

func TestState_BitFlags(t *testing.T) {
	// Test that states can be combined as bit flags
	combined := Continue | HasChildren
	if combined&Continue == 0 {
		t.Error("Combined state should include Continue")
	}
	if combined&HasChildren == 0 {
		t.Error("Combined state should include HasChildren")
	}
	if combined&Close != 0 {
		t.Error("Combined state should not include Close")
	}
}

func TestNewContextKey(t *testing.T) {
	key1 := NewContextKey()
	key2 := NewContextKey()

	if key1 == key2 {
		t.Error("NewContextKey should return unique keys")
	}
	if key2 <= key1 {
		t.Error("NewContextKey should return incrementing values")
	}
}

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg.Options == nil {
		t.Error("Config.Options should not be nil")
	}
	if cfg.BlockParsers == nil {
		t.Error("Config.BlockParsers should not be nil")
	}
	if cfg.InlineParsers == nil {
		t.Error("Config.InlineParsers should not be nil")
	}
	if cfg.ParagraphTransformers == nil {
		t.Error("Config.ParagraphTransformers should not be nil")
	}
	if cfg.ASTTransformers == nil {
		t.Error("Config.ASTTransformers should not be nil")
	}
}

func TestWithAttribute(t *testing.T) {
	opt := WithAttribute()
	cfg := NewConfig()

	opt.SetParserOption(cfg)

	if val, ok := cfg.Options[optAttribute]; !ok || val != true {
		t.Error("WithAttribute should set optAttribute to true")
	}
}

func TestWithEscapedSpace(t *testing.T) {
	opt := WithEscapedSpace()
	cfg := NewConfig()

	opt.SetParserOption(cfg)

	if !cfg.EscapedSpace {
		t.Error("WithEscapedSpace should set EscapedSpace to true")
	}
}

func TestWithOption(t *testing.T) {
	testName := OptionName("test")
	testValue := "value"

	opt := WithOption(testName, testValue)
	cfg := NewConfig()

	opt.SetParserOption(cfg)

	if val, ok := cfg.Options[testName]; !ok || val != testValue {
		t.Errorf("WithOption should set option %q to %v", testName, testValue)
	}
}

func TestDefaultBlockParsers(t *testing.T) {
	parsers := DefaultBlockParsers()

	if len(parsers) == 0 {
		t.Error("DefaultBlockParsers should return non-empty slice")
	}

	// Check that all parsers are actually BlockParsers
	for _, p := range parsers {
		if _, ok := p.Value.(BlockParser); !ok {
			t.Errorf("Parser %v is not a BlockParser", p.Value)
		}
	}
}

func TestDefaultInlineParsers(t *testing.T) {
	parsers := DefaultInlineParsers()

	if len(parsers) == 0 {
		t.Error("DefaultInlineParsers should return non-empty slice")
	}

	// Check that all parsers are actually InlineParsers
	for _, p := range parsers {
		if _, ok := p.Value.(InlineParser); !ok {
			t.Errorf("Parser %v is not an InlineParser", p.Value)
		}
	}
}

func TestDefaultParagraphTransformers(t *testing.T) {
	transformers := DefaultParagraphTransformers()

	if len(transformers) == 0 {
		t.Error("DefaultParagraphTransformers should return non-empty slice")
	}

	// Check that all transformers are actually ParagraphTransformers
	for _, tr := range transformers {
		if _, ok := tr.Value.(ParagraphTransformer); !ok {
			t.Errorf("Transformer %v is not a ParagraphTransformer", tr.Value)
		}
	}
}
