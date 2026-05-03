package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestLinkParser_Trigger(t *testing.T) {
	parser := NewLinkParser()
	triggers := parser.Trigger()

	expected := []byte{'!', '[', ']'}
	if len(triggers) != len(expected) {
		t.Fatalf("Trigger() length = %d, want %d", len(triggers), len(expected))
	}

	for i, trigger := range expected {
		if triggers[i] != trigger {
			t.Errorf("Trigger()[%d] = %c, want %c", i, triggers[i], trigger)
		}
	}
}

func TestLinkParser_Parse_OpenBracket(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string
	}{
		{
			name:     "opening bracket creates state",
			input:    "[text",
			wantType: "*parser.linkLabelState",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewLinkParser()
			source := []byte(tc.input)
			reader := text.NewReader(source)
			pc := NewContext()
			parent := ast.NewDocument()

			// Parse opening bracket
			node := parser.Parse(parent, reader, pc)

			if node == nil {
				t.Fatal("Parse() returned nil for opening bracket")
			}

			state, ok := node.(*linkLabelState)
			if !ok {
				t.Fatalf("Parse() returned %T, want *linkLabelState", node)
			}

			if state.IsImage {
				t.Error("linkLabelState.IsImage = true, want false for regular link")
			}
		})
	}
}

func TestLinkParser_Parse_Image(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantDest string
		wantAlt  string
	}{
		{
			name:     "simple image",
			input:    "![alt](image.png)",
			wantDest: "image.png",
			wantAlt:  "alt",
		},
		{
			name:     "image with empty alt",
			input:    "![](image.png)",
			wantDest: "image.png",
			wantAlt:  "",
		},
		{
			name:     "image with path",
			input:    "![photo](/images/photo.jpg)",
			wantDest: "/images/photo.jpg",
			wantAlt:  "photo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewLinkParser()
			source := []byte(tc.input)
			reader := text.NewReader(source)
			pc := NewContext()
			parent := ast.NewDocument()

			// Parse opening "!["
			node := parser.Parse(parent, reader, pc)

			if node == nil {
				t.Fatal("Parse() returned nil for image opener")
			}

			state, ok := node.(*linkLabelState)
			if !ok {
				t.Fatalf("Parse() for image opener returned %T, want *linkLabelState", node)
			}

			if !state.IsImage {
				t.Error("linkLabelState.IsImage = false, want true")
			}
		})
	}
}

func TestLinkParser_ContainsLink(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() ast.Node
		expected bool
	}{
		{
			name: "nil node",
			setup: func() ast.Node {
				return nil
			},
			expected: false,
		},
		{
			name: "node with link child",
			setup: func() ast.Node {
				parent := ast.NewDocument()
				link := ast.NewLink()
				parent.AppendChild(parent, link)
				return parent
			},
			expected: true,
		},
		{
			name: "node with text child only",
			setup: func() ast.Node {
				parent := ast.NewDocument()
				text := ast.NewString([]byte("text"))
				parent.AppendChild(parent, text)
				return parent
			},
			expected: false,
		},
		{
			name: "nested link in grandchild",
			setup: func() ast.Node {
				parent := ast.NewDocument()
				child := ast.NewEmphasis(1)
				grandchild := ast.NewLink()
				parent.AppendChild(parent, child)
				child.AppendChild(child, grandchild)
				return parent
			},
			expected: true,
		},
		{
			name: "link as sibling",
			setup: func() ast.Node {
				parent := ast.NewDocument()
				text := ast.NewString([]byte("text"))
				link := ast.NewLink()
				parent.AppendChild(parent, text)
				parent.AppendChild(parent, link)
				return text
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &linkParser{}
			node := tc.setup()
			result := parser.containsLink(node)

			if result != tc.expected {
				t.Errorf("containsLink() = %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestLinkLabelState_Text(t *testing.T) {
	tests := []struct {
		name     string
		segment  text.Segment
		source   []byte
		expected string
	}{
		{
			name:     "extract simple text",
			segment:  text.NewSegment(0, 4),
			source:   []byte("test content"),
			expected: "test",
		},
		{
			name:     "extract middle segment",
			segment:  text.NewSegment(5, 12),
			source:   []byte("test content here"),
			expected: "content",
		},
		{
			name:     "empty segment",
			segment:  text.NewSegment(0, 0),
			source:   []byte("test"),
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := &linkLabelState{
				Segment: tc.segment,
			}

			result := string(state.Text(tc.source))
			if result != tc.expected {
				t.Errorf("Text() = %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestLinkLabelState_Kind(t *testing.T) {
	state := &linkLabelState{}
	kind := state.Kind()

	if kind.String() != "LinkLabelState" {
		t.Errorf("Kind().String() = %q, want %q", kind.String(), "LinkLabelState")
	}
}

func TestNewLinkLabelState(t *testing.T) {
	tests := []struct {
		name    string
		segment text.Segment
		isImage bool
	}{
		{
			name:    "create link state",
			segment: text.NewSegment(0, 5),
			isImage: false,
		},
		{
			name:    "create image state",
			segment: text.NewSegment(0, 6),
			isImage: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := newLinkLabelState(tc.segment, tc.isImage)

			if state == nil {
				t.Fatal("newLinkLabelState() returned nil")
			}

			if state.Segment != tc.segment {
				t.Errorf("Segment = %v, want %v", state.Segment, tc.segment)
			}

			if state.IsImage != tc.isImage {
				t.Errorf("IsImage = %v, want %v", state.IsImage, tc.isImage)
			}
		})
	}
}

func TestLinkLabelStateLength(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *linkLabelState
		expected int
	}{
		{
			name: "nil state",
			setup: func() *linkLabelState {
				return nil
			},
			expected: 0,
		},
		{
			name: "state with nil Last",
			setup: func() *linkLabelState {
				state := &linkLabelState{}
				state.First = state
				return state
			},
			expected: 0,
		},
		{
			name: "state with nil First",
			setup: func() *linkLabelState {
				state := &linkLabelState{}
				state.Last = state
				return state
			},
			expected: 0,
		},
		{
			name: "single state",
			setup: func() *linkLabelState {
				state := &linkLabelState{
					Segment: text.NewSegment(0, 10),
				}
				state.First = state
				state.Last = state
				return state
			},
			expected: 10,
		},
		{
			name: "multiple states",
			setup: func() *linkLabelState {
				first := &linkLabelState{
					Segment: text.NewSegment(0, 5),
				}
				last := &linkLabelState{
					Segment: text.NewSegment(5, 15),
				}
				first.First = first
				first.Last = last
				last.First = first
				last.Last = last
				return first
			},
			expected: 15,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			state := tc.setup()
			result := linkLabelStateLength(state)

			if result != tc.expected {
				t.Errorf("linkLabelStateLength() = %d, want %d", result, tc.expected)
			}
		})
	}
}

func TestPushLinkLabelState(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (Context, *linkLabelState)
		check func(*testing.T, Context)
	}{
		{
			name: "push to empty context",
			setup: func() (Context, *linkLabelState) {
				pc := NewContext()
				state := newLinkLabelState(text.NewSegment(0, 5), false)
				return pc, state
			},
			check: func(t *testing.T, pc Context) {
				val := pc.Get(linkLabelStateKey)
				if val == nil {
					t.Fatal("context key not set")
				}
				state := val.(*linkLabelState)
				if state.First != state {
					t.Error("First not set to self")
				}
				if state.Last != state {
					t.Error("Last not set to self")
				}
			},
		},
		{
			name: "push second state",
			setup: func() (Context, *linkLabelState) {
				pc := NewContext()
				first := newLinkLabelState(text.NewSegment(0, 5), false)
				pushLinkLabelState(pc, first)
				second := newLinkLabelState(text.NewSegment(5, 10), false)
				return pc, second
			},
			check: func(t *testing.T, pc Context) {
				val := pc.Get(linkLabelStateKey)
				if val == nil {
					t.Fatal("context key not set")
				}
				state := val.(*linkLabelState)
				if state.Last == state {
					t.Error("Last should point to second state")
				}
				if state.Last.Prev == nil {
					t.Error("Last.Prev should point to first")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pc, state := tc.setup()
			pushLinkLabelState(pc, state)
			tc.check(t, pc)
		})
	}
}

func TestRemoveLinkLabelState(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (Context, *linkLabelState)
		check func(*testing.T, Context)
	}{
		{
			name: "remove from empty context",
			setup: func() (Context, *linkLabelState) {
				pc := NewContext()
				state := newLinkLabelState(text.NewSegment(0, 5), false)
				return pc, state
			},
			check: func(t *testing.T, pc Context) {
				// Should not panic, context remains empty
				val := pc.Get(linkLabelStateKey)
				if val != nil {
					t.Error("context should remain empty")
				}
			},
		},
		{
			name: "remove only state",
			setup: func() (Context, *linkLabelState) {
				pc := NewContext()
				state := newLinkLabelState(text.NewSegment(0, 5), false)
				pushLinkLabelState(pc, state)
				return pc, state
			},
			check: func(t *testing.T, pc Context) {
				val := pc.Get(linkLabelStateKey)
				if val != nil {
					t.Error("context should be cleared after removing only state")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pc, state := tc.setup()
			removeLinkLabelState(pc, state)
			tc.check(t, pc)
		})
	}
}

func TestProcessLinkLabelOpen(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		pos       int
		isImage   bool
		wantStart int
		wantStop  int
	}{
		{
			name:      "normal link opening",
			input:     "[text]",
			pos:       0,
			isImage:   false,
			wantStart: 0,
			wantStop:  1,
		},
		{
			name:      "image opening",
			input:     "![alt]",
			pos:       1,
			isImage:   true,
			wantStart: 0,
			wantStop:  2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

			state := processLinkLabelOpen(reader, tc.pos, tc.isImage, pc)

			if state == nil {
				t.Fatal("processLinkLabelOpen() returned nil")
			}

			if state.Segment.Start != tc.wantStart {
				t.Errorf("Segment.Start = %d, want %d", state.Segment.Start, tc.wantStart)
			}

			if state.Segment.Stop != tc.wantStop {
				t.Errorf("Segment.Stop = %d, want %d", state.Segment.Stop, tc.wantStop)
			}

			if state.IsImage != tc.isImage {
				t.Errorf("IsImage = %v, want %v", state.IsImage, tc.isImage)
			}

			// Verify state was pushed to context
			val := pc.Get(linkLabelStateKey)
			if val == nil {
				t.Error("state not pushed to context")
			}
		})
	}
}
