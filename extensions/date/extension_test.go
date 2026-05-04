package date

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/stretchr/testify/assert"
)

const (
	calendarLinkName = "Calendar"
)

func TestDateExtensionName(t *testing.T) {
	ext := Date{}
	assert.Equal(t, "date", ext.Name())
}

func TestCalendarCommand(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedIcon string
		expectedName string
		expectedHref string
	}{
		{
			name:         "Calendar icon",
			method:       "Icon",
			expectedIcon: "fa-regular fa-calendar-days",
		},
		{
			name:         "Calendar name",
			method:       "Name",
			expectedName: calendarLinkName,
		},
		{
			name:         "Calendar href",
			method:       "Attrs",
			expectedHref: "/+/calendar",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cal := Calendar{}

			switch tc.method {
			case "Icon":
				assert.Equal(t, tc.expectedIcon, cal.Icon())
			case "Name":
				assert.Equal(t, tc.expectedName, cal.Name())
			case "Attrs":
				attrs := cal.Attrs()
				href, ok := attrs["href"]
				assert.True(t, ok, "href attribute should exist")
				assert.Equal(t, tc.expectedHref, href)
			}
		})
	}
}

func TestLinks(t *testing.T) {
	type mockPage struct {
		xlog.Page
	}

	tests := []struct {
		name          string
		expectedCount int
	}{
		{
			name:          "returns calendar command",
			expectedCount: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var mock mockPage
			commands := links(&mock)
			assert.Len(t, commands, tc.expectedCount)
			assert.IsType(t, Calendar{}, commands[0])
		})
	}
}

func TestDateNode(t *testing.T) {
	tests := []struct {
		name   string
		method string
	}{
		{
			name:   "Kind returns correct NodeKind",
			method: "Kind",
		},
		{
			name:   "Dump outputs node information",
			method: "Dump",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := &DateNode{
				time: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			}

			switch tc.method {
			case "Kind":
				kind := node.Kind()
				assert.Equal(t, KindDate, kind)
			case "Dump":
				// Dump should not panic
				assert.NotPanics(t, func() {
					node.Dump([]byte("test source"), 0)
				})
			}
		})
	}
}

func TestDateRenderer_Render(t *testing.T) {
	tests := []struct {
		name        string
		node        ast.Node
		entering    bool
		expectWrite bool
		expectHref  string
	}{
		{
			name: "renders date node on entering",
			node: &DateNode{
				time: time.Date(2023, time.May, 15, 0, 0, 0, 0, time.UTC),
			},
			entering:    true,
			expectWrite: true,
			expectHref:  "/+/date/15-5-2023",
		},
		{
			name: "skips on not entering",
			node: &DateNode{
				time: time.Date(2023, time.May, 15, 0, 0, 0, 0, time.UTC),
			},
			entering:    false,
			expectWrite: false,
		},
		{
			name: "skips non-DateNode",
			node: &ast.Text{
				BaseInline: ast.BaseInline{},
			},
			entering:    true,
			expectWrite: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)
			renderer := &dateRenderer{}

			status, err := renderer.render(writer, []byte{}, tc.node, tc.entering)

			assert.NoError(t, err)
			assert.Equal(t, ast.WalkContinue, status)

			// Flush to get the written content
			if err := writer.Flush(); err != nil {
				t.Fatalf("Failed to flush buffer: %v", err)
			}

			if tc.expectWrite {
				output := buf.String()
				assert.Contains(t, output, tc.expectHref, "Output should contain href")
				assert.Contains(t, output, "15 May 2023", "Output should contain formatted date")
				assert.Contains(t, output, "fa-regular fa-clock", "Output should contain clock icon")
			} else {
				assert.Empty(t, buf.String(), "Should not write anything")
			}
		})
	}
}
