package date

import (
	"testing"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestDateParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantDate string // Expected date in YYYY-MM-DD format
		wantNil  bool
	}{
		{
			name:     "ISO format with dashes",
			input:    " 2026-3-26",
			wantDate: "2026-03-26",
		},
		{
			name:     "Full month name with dashes",
			input:    " 2026-March-26",
			wantDate: "2026-03-26",
		},
		{
			name:     "Full month name with slashes",
			input:    " 2026/March/26",
			wantDate: "2026-03-26",
		},
		{
			name:     "Full month name with backslashes",
			input:    " 2026\\March\\26",
			wantDate: "2026-03-26",
		},
		{
			name:     "Short month name with dashes",
			input:    " 2026-Mar-26",
			wantDate: "2026-03-26",
		},
		{
			name:     "Day first format with full month",
			input:    " 26-March-2026",
			wantDate: "2026-03-26",
		},
		{
			name:     "Day first format with short month",
			input:    " 26/Mar/2026",
			wantDate: "2026-03-26",
		},
		{
			name:     "Month first format with short month",
			input:    " Mar-26-2026",
			wantDate: "2026-03-26",
		},
		{
			name:     "Month first format with full month",
			input:    " March/26/2026",
			wantDate: "2026-03-26",
		},
		{
			name:    "Invalid date format",
			input:   " not-a-date",
			wantNil: true,
		},
		{
			name:     "Text after date is ignored",
			input:    " 2026-03-26-extra",
			wantDate: "2026-03-26",
		},
		{
			name:    "Empty input",
			input:   "",
			wantNil: true,
		},
		{
			name:     "Date without space prefix still parses",
			input:    "2026-3-26",
			wantDate: "2026-03-26",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &dateParser{}
			reader := text.NewReader([]byte(tt.input))
			pc := parser.NewContext()

			node := p.Parse(ast.NewDocument(), reader, pc)

			if tt.wantNil {
				if node != nil {
					t.Errorf("Parse() expected nil, got node")
				}
				return
			}

			if node == nil {
				t.Fatalf("Parse() returned nil, expected DateNode")
			}

			dateNode, ok := node.(*DateNode)
			if !ok {
				t.Fatalf("Parse() returned %T, expected *DateNode", node)
			}

			gotDate := dateNode.time.Format("2006-01-02")
			if gotDate != tt.wantDate {
				t.Errorf("Parse() date = %v, want %v", gotDate, tt.wantDate)
			}
		})
	}
}

func TestDateParser_Trigger(t *testing.T) {
	p := &dateParser{}
	trigger := p.Trigger()

	if len(trigger) != 1 {
		t.Fatalf("Trigger() length = %d, want 1", len(trigger))
	}

	if trigger[0] != ' ' {
		t.Errorf("Trigger() = %q, want %q", trigger[0], ' ')
	}
}

func TestDateParser_ParseVariousYears(t *testing.T) {
	tests := []struct {
		year  int
		month time.Month
		day   int
	}{
		{2020, time.January, 1},
		{2025, time.December, 31},
		{2026, time.March, 26},
		{1999, time.June, 15},
	}

	p := &dateParser{}

	for _, tt := range tests {
		input := " " + time.Date(tt.year, tt.month, tt.day, 0, 0, 0, 0, time.UTC).Format("2006-1-2")
		reader := text.NewReader([]byte(input))
		pc := parser.NewContext()

		node := p.Parse(ast.NewDocument(), reader, pc)
		if node == nil {
			t.Fatalf("Parse(%s) returned nil", input)
		}

		dateNode := node.(*DateNode)
		if dateNode.time.Year() != tt.year || dateNode.time.Month() != tt.month || dateNode.time.Day() != tt.day {
			t.Errorf("Parse(%s) = %v, want %04d-%02d-%02d", input, dateNode.time, tt.year, tt.month, tt.day)
		}
	}
}
