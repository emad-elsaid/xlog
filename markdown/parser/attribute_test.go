package parser

import (
	"bytes"
	"testing"

	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestAttributes_Find(t *testing.T) {
	tests := []struct {
		name       string
		attrs      Attributes
		searchName []byte
		wantValue  any
		wantFound  bool
	}{
		{
			name: "find existing id attribute",
			attrs: Attributes{
				{Name: []byte("id"), Value: []byte("test-id")},
				{Name: []byte("class"), Value: []byte("container")},
			},
			searchName: []byte("id"),
			wantValue:  []byte("test-id"),
			wantFound:  true,
		},
		{
			name: "find existing class attribute",
			attrs: Attributes{
				{Name: []byte("id"), Value: []byte("test-id")},
				{Name: []byte("class"), Value: []byte("container")},
			},
			searchName: []byte("class"),
			wantValue:  []byte("container"),
			wantFound:  true,
		},
		{
			name: "find non-existent attribute",
			attrs: Attributes{
				{Name: []byte("id"), Value: []byte("test-id")},
			},
			searchName: []byte("class"),
			wantValue:  nil,
			wantFound:  false,
		},
		{
			name:       "search in empty attributes",
			attrs:      Attributes{},
			searchName: []byte("id"),
			wantValue:  nil,
			wantFound:  false,
		},
		{
			name: "find attribute with number value",
			attrs: Attributes{
				{Name: []byte("width"), Value: 100.0},
			},
			searchName: []byte("width"),
			wantValue:  100.0,
			wantFound:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotValue, gotFound := tc.attrs.Find(tc.searchName)
			if gotFound != tc.wantFound {
				t.Errorf("Find() found = %v, want %v", gotFound, tc.wantFound)
			}
			if tc.wantFound {
				switch want := tc.wantValue.(type) {
				case []byte:
					if got, ok := gotValue.([]byte); !ok || !bytes.Equal(got, want) {
						t.Errorf("Find() value = %v, want %v", gotValue, tc.wantValue)
					}
				default:
					if gotValue != tc.wantValue {
						t.Errorf("Find() value = %v, want %v", gotValue, tc.wantValue)
					}
				}
			}
		})
	}
}

func TestAttributes_findUpdate(t *testing.T) {
	tests := []struct {
		name       string
		attrs      Attributes
		searchName []byte
		updateFunc func(v any) any
		wantFound  bool
		wantValue  any
	}{
		{
			name: "update existing attribute",
			attrs: Attributes{
				{Name: []byte("class"), Value: []byte("old")},
			},
			searchName: []byte("class"),
			updateFunc: func(v any) any {
				return append(v.([]byte), []byte(" new")...)
			},
			wantFound: true,
			wantValue: []byte("old new"),
		},
		{
			name: "update non-existent attribute",
			attrs: Attributes{
				{Name: []byte("id"), Value: []byte("test")},
			},
			searchName: []byte("class"),
			updateFunc: func(v any) any {
				return []byte("new")
			},
			wantFound: false,
			wantValue: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotFound := tc.attrs.findUpdate(tc.searchName, tc.updateFunc)
			if gotFound != tc.wantFound {
				t.Errorf("findUpdate() found = %v, want %v", gotFound, tc.wantFound)
			}
			if tc.wantFound {
				val, _ := tc.attrs.Find(tc.searchName)
				if !bytes.Equal(val.([]byte), tc.wantValue.([]byte)) {
					t.Errorf("Updated value = %v, want %v", val, tc.wantValue)
				}
			}
		})
	}
}

func TestParseAttributes(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantAttrs Attributes
		wantOk    bool
	}{
		{
			name:  "simple id shorthand",
			input: "{#myid}",
			wantAttrs: Attributes{
				{Name: []byte("id"), Value: []byte("myid")},
			},
			wantOk: true,
		},
		{
			name:  "simple class shorthand",
			input: "{.myclass}",
			wantAttrs: Attributes{
				{Name: []byte("class"), Value: []byte("myclass")},
			},
			wantOk: true,
		},
		{
			name:  "id and class",
			input: "{#myid .myclass}",
			wantAttrs: Attributes{
				{Name: []byte("id"), Value: []byte("myid")},
				{Name: []byte("class"), Value: []byte("myclass")},
			},
			wantOk: true,
		},
		{
			name:  "multiple classes concatenated",
			input: "{.class1 .class2}",
			wantAttrs: Attributes{
				{Name: []byte("class"), Value: []byte("class1 class2")},
			},
			wantOk: true,
		},
		{
			name:  "name-value string attribute",
			input: `{title="Hello World"}`,
			wantAttrs: Attributes{
				{Name: []byte("title"), Value: []byte("Hello World")},
			},
			wantOk: true,
		},
		{
			name:  "name-value number attribute",
			input: "{width=100}",
			wantAttrs: Attributes{
				{Name: []byte("width"), Value: 100.0},
			},
			wantOk: true,
		},
		{
			name:  "name-value negative number",
			input: "{offset=-50}",
			wantAttrs: Attributes{
				{Name: []byte("offset"), Value: -50.0},
			},
			wantOk: true,
		},
		{
			name:  "name-value float",
			input: "{ratio=1.5}",
			wantAttrs: Attributes{
				{Name: []byte("ratio"), Value: 1.5},
			},
			wantOk: true,
		},
		{
			name:  "name-value boolean true",
			input: "{enabled=true}",
			wantAttrs: Attributes{
				{Name: []byte("enabled"), Value: true},
			},
			wantOk: true,
		},
		{
			name:  "name-value boolean false",
			input: "{disabled=false}",
			wantAttrs: Attributes{
				{Name: []byte("disabled"), Value: false},
			},
			wantOk: true,
		},
		{
			name:  "name-value null",
			input: "{value=null}",
			wantAttrs: Attributes{
				{Name: []byte("value"), Value: nil},
			},
			wantOk: true,
		},
		{
			name:  "array attribute",
			input: `{items=[1, 2, 3]}`,
			wantAttrs: Attributes{
				{Name: []byte("items"), Value: []any{1.0, 2.0, 3.0}},
			},
			wantOk: true,
		},
		{
			name:  "mixed array",
			input: `{mixed=[1, "text", true]}`,
			wantAttrs: Attributes{
				{Name: []byte("mixed"), Value: []any{1.0, []byte("text"), true}},
			},
			wantOk: true,
		},
		{
			name:  "empty array",
			input: "{items=[]}",
			wantAttrs: Attributes{
				{Name: []byte("items"), Value: []any{}},
			},
			wantOk: true,
		},
		{
			name:  "escaped characters in string",
			input: `{text="line1\nline2\ttab"}`,
			wantAttrs: Attributes{
				{Name: []byte("text"), Value: []byte("line1\nline2\ttab")},
			},
			wantOk: true,
		},
		{
			name:  "escaped quote in string",
			input: `{quote="He said \"Hello\""}`,
			wantAttrs: Attributes{
				{Name: []byte("quote"), Value: []byte("He said \"Hello\"")},
			},
			wantOk: true,
		},
		{
			name:  "comma separated attributes",
			input: "{width=100, height=200}",
			wantAttrs: Attributes{
				{Name: []byte("width"), Value: 100.0},
				{Name: []byte("height"), Value: 200.0},
			},
			wantOk: true,
		},
		{
			name:  "complex mixed attributes",
			input: `{#id .class1 .class2, title="Test", width=100}`,
			wantAttrs: Attributes{
				{Name: []byte("id"), Value: []byte("id")},
				{Name: []byte("class"), Value: []byte("class1 class2")},
				{Name: []byte("title"), Value: []byte("Test")},
				{Name: []byte("width"), Value: 100.0},
			},
			wantOk: true,
		},
		{
			name:   "no opening brace",
			input:  "id=test}",
			wantOk: false,
		},
		{
			name:   "no closing brace",
			input:  "{id=test",
			wantOk: false,
		},
		{
			name:   "empty braces",
			input:  "{}",
			wantOk: true,
		},
		{
			name:  "whitespace handling",
			input: "{  #id   .class  }",
			wantAttrs: Attributes{
				{Name: []byte("id"), Value: []byte("id")},
				{Name: []byte("class"), Value: []byte("class")},
			},
			wantOk: true,
		},
		{
			name:  "scientific notation",
			input: "{value=1.5e10}",
			wantAttrs: Attributes{
				{Name: []byte("value"), Value: 1.5e10},
			},
			wantOk: true,
		},
		{
			name:  "negative scientific notation",
			input: "{value=-2.5e-3}",
			wantAttrs: Attributes{
				{Name: []byte("value"), Value: -2.5e-3},
			},
			wantOk: true,
		},
		{
			name:  "id with hyphens and underscores",
			input: "{#my-id_123}",
			wantAttrs: Attributes{
				{Name: []byte("id"), Value: []byte("my-id_123")},
			},
			wantOk: true,
		},
		{
			name:  "class with dots (CSS module naming)",
			input: "{.component.active}",
			wantAttrs: Attributes{
				{Name: []byte("class"), Value: []byte("component.active")},
			},
			wantOk: true,
		},
		{
			name:  "attribute name with colon (data attributes)",
			input: "{data:value=123}",
			wantAttrs: Attributes{
				{Name: []byte("data:value"), Value: 123.0},
			},
			wantOk: true,
		},
		{
			name:  "nested attributes",
			input: `{meta={author="John", year=2024}}`,
			wantAttrs: Attributes{
				{
					Name: []byte("meta"),
					Value: Attributes{
						{Name: []byte("author"), Value: []byte("John")},
						{Name: []byte("year"), Value: 2024.0},
					},
				},
			},
			wantOk: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.input))
			attrs, ok := ParseAttributes(reader)

			if ok != tc.wantOk {
				t.Errorf("ParseAttributes() ok = %v, want %v", ok, tc.wantOk)
				return
			}

			if !tc.wantOk {
				return
			}

			if len(attrs) != len(tc.wantAttrs) {
				t.Errorf("ParseAttributes() got %d attributes, want %d", len(attrs), len(tc.wantAttrs))
				t.Logf("Got: %+v", attrs)
				t.Logf("Want: %+v", tc.wantAttrs)
				return
			}

			for i, wantAttr := range tc.wantAttrs {
				gotAttr := attrs[i]
				if !bytes.Equal(gotAttr.Name, wantAttr.Name) {
					t.Errorf("Attribute[%d] name = %s, want %s", i, gotAttr.Name, wantAttr.Name)
				}

				compareAttributeValue(t, i, gotAttr.Value, wantAttr.Value)
			}
		})
	}
}

func compareAttributeValue(t *testing.T, index int, got, want any) {
	t.Helper()

	switch wantVal := want.(type) {
	case []byte:
		gotVal, ok := got.([]byte)
		if !ok {
			t.Errorf("Attribute[%d] value type = %T, want []byte", index, got)
			return
		}
		if !bytes.Equal(gotVal, wantVal) {
			t.Errorf("Attribute[%d] value = %q, want %q", index, gotVal, wantVal)
		}
	case []any:
		gotVal, ok := got.([]any)
		if !ok {
			t.Errorf("Attribute[%d] value type = %T, want []any", index, got)
			return
		}
		if len(gotVal) != len(wantVal) {
			t.Errorf("Attribute[%d] array length = %d, want %d", index, len(gotVal), len(wantVal))
			return
		}
		for j := range wantVal {
			compareAttributeValue(t, index, gotVal[j], wantVal[j])
		}
	case Attributes:
		gotVal, ok := got.(Attributes)
		if !ok {
			t.Errorf("Attribute[%d] value type = %T, want Attributes", index, got)
			return
		}
		if len(gotVal) != len(wantVal) {
			t.Errorf("Attribute[%d] nested attributes length = %d, want %d", index, len(gotVal), len(wantVal))
			return
		}
		for j, wantNestedAttr := range wantVal {
			gotNestedAttr := gotVal[j]
			if !bytes.Equal(gotNestedAttr.Name, wantNestedAttr.Name) {
				t.Errorf("Nested attribute[%d][%d] name = %s, want %s",
					index, j, gotNestedAttr.Name, wantNestedAttr.Name)
			}
			compareAttributeValue(t, j, gotNestedAttr.Value, wantNestedAttr.Value)
		}
	default:
		if got != want {
			t.Errorf("Attribute[%d] value = %v (%T), want %v (%T)", index, got, got, want, want)
		}
	}
}

func TestParseAttributeString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   []byte
		wantOk bool
	}{
		{
			name:   "simple string",
			input:  `"hello"`,
			want:   []byte("hello"),
			wantOk: true,
		},
		{
			name:   "empty string",
			input:  `""`,
			want:   []byte(""),
			wantOk: true,
		},
		{
			name:   "string with spaces",
			input:  `"hello world"`,
			want:   []byte("hello world"),
			wantOk: true,
		},
		{
			name:   "escaped quote",
			input:  `"say \"hi\""`,
			want:   []byte(`say "hi"`),
			wantOk: true,
		},
		{
			name:   "escaped backslash",
			input:  `"path\\to\\file"`,
			want:   []byte(`path\to\file`),
			wantOk: true,
		},
		{
			name:   "escaped newline",
			input:  `"line1\nline2"`,
			want:   []byte("line1\nline2"),
			wantOk: true,
		},
		{
			name:   "escaped tab",
			input:  `"col1\tcol2"`,
			want:   []byte("col1\tcol2"),
			wantOk: true,
		},
		{
			name:   "escaped carriage return",
			input:  `"text\r"`,
			want:   []byte("text\r"),
			wantOk: true,
		},
		{
			name:   "escaped formfeed",
			input:  `"text\f"`,
			want:   []byte("text\f"),
			wantOk: true,
		},
		{
			name:   "escaped backspace",
			input:  `"text\b"`,
			want:   []byte("text\b"),
			wantOk: true,
		},
		{
			name:   "escaped slash",
			input:  `"url\/path"`,
			want:   []byte("url/path"),
			wantOk: true,
		},
		{
			name:   "unknown escape sequence",
			input:  `"test\x"`,
			want:   []byte(`test\x`),
			wantOk: true,
		},
		{
			name:   "unterminated string",
			input:  `"hello`,
			wantOk: false,
		},
		{
			name:   "backslash at end",
			input:  `"test\"`,
			wantOk: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.input))
			got, ok := parseAttributeString(reader)

			if ok != tc.wantOk {
				t.Errorf("parseAttributeString() ok = %v, want %v", ok, tc.wantOk)
				return
			}

			if tc.wantOk && !bytes.Equal(got, tc.want) {
				t.Errorf("parseAttributeString() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseAttributeNumber(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   float64
		wantOk bool
	}{
		{
			name:   "positive integer",
			input:  "42",
			want:   42.0,
			wantOk: true,
		},
		{
			name:   "negative integer",
			input:  "-42",
			want:   -42.0,
			wantOk: true,
		},
		{
			name:   "positive sign",
			input:  "+42",
			want:   42.0,
			wantOk: true,
		},
		{
			name:   "zero",
			input:  "0",
			want:   0.0,
			wantOk: true,
		},
		{
			name:   "float",
			input:  "3.14",
			want:   3.14,
			wantOk: true,
		},
		{
			name:   "negative float",
			input:  "-2.5",
			want:   -2.5,
			wantOk: true,
		},
		{
			name:   "scientific notation",
			input:  "1.5e10",
			want:   1.5e10,
			wantOk: true,
		},
		{
			name:   "negative exponent",
			input:  "2.5e-3",
			want:   2.5e-3,
			wantOk: true,
		},
		{
			name:   "uppercase E",
			input:  "1.5E10",
			want:   1.5e10,
			wantOk: true,
		},
		{
			name:   "positive exponent sign",
			input:  "1e+5",
			want:   1e+5,
			wantOk: true,
		},
		{
			name:   "leading zeros",
			input:  "007",
			want:   7.0,
			wantOk: true,
		},
		{
			name:   "decimal without integer part",
			input:  ".5",
			wantOk: false,
		},
		{
			name:   "just a sign",
			input:  "-",
			wantOk: false,
		},
		{
			name:   "just a plus",
			input:  "+",
			wantOk: false,
		},
		{
			name:   "not a number",
			input:  "abc",
			wantOk: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.input))
			got, ok := parseAttributeNumber(reader)

			if ok != tc.wantOk {
				t.Errorf("parseAttributeNumber() ok = %v, want %v", ok, tc.wantOk)
				return
			}

			if tc.wantOk && got != tc.want {
				t.Errorf("parseAttributeNumber() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestParseAttributeOthers(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   any
		wantOk bool
	}{
		{
			name:   "boolean true",
			input:  "true",
			want:   true,
			wantOk: true,
		},
		{
			name:   "boolean false",
			input:  "false",
			want:   false,
			wantOk: true,
		},
		{
			name:   "null value",
			input:  "null",
			want:   nil,
			wantOk: true,
		},
		{
			name:   "identifier with letters",
			input:  "myValue",
			want:   []byte("myValue"),
			wantOk: true,
		},
		{
			name:   "identifier with underscores",
			input:  "my_value",
			want:   []byte("my_value"),
			wantOk: true,
		},
		{
			name:   "identifier with hyphens",
			input:  "my-value",
			want:   []byte("my-value"),
			wantOk: true,
		},
		{
			name:   "identifier with colons",
			input:  "data:value",
			want:   []byte("data:value"),
			wantOk: true,
		},
		{
			name:   "identifier with dots",
			input:  "obj.prop",
			want:   []byte("obj.prop"),
			wantOk: true,
		},
		{
			name:   "identifier with numbers",
			input:  "value123",
			want:   []byte("value123"),
			wantOk: true,
		},
		{
			name:   "uppercase identifier",
			input:  "VALUE",
			want:   []byte("VALUE"),
			wantOk: true,
		},
		{
			name:   "starting with number",
			input:  "123abc",
			wantOk: false,
		},
		{
			name:   "starting with hyphen",
			input:  "-value",
			wantOk: false,
		},
		{
			name:   "starting with dot",
			input:  ".value",
			wantOk: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.input))
			got, ok := parseAttributeOthers(reader)

			if ok != tc.wantOk {
				t.Errorf("parseAttributeOthers() ok = %v, want %v", ok, tc.wantOk)
				return
			}

			if !tc.wantOk {
				return
			}

			switch wantVal := tc.want.(type) {
			case []byte:
				gotVal, ok := got.([]byte)
				if !ok {
					t.Errorf("parseAttributeOthers() type = %T, want []byte", got)
					return
				}
				if !bytes.Equal(gotVal, wantVal) {
					t.Errorf("parseAttributeOthers() = %q, want %q", gotVal, wantVal)
				}
			default:
				if got != tc.want {
					t.Errorf("parseAttributeOthers() = %v, want %v", got, tc.want)
				}
			}
		})
	}
}

func TestParseAttributeArray(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   []any
		wantOk bool
	}{
		{
			name:   "empty array",
			input:  "[]",
			want:   []any{},
			wantOk: true,
		},
		{
			name:   "single number",
			input:  "[42]",
			want:   []any{42.0},
			wantOk: true,
		},
		{
			name:   "multiple numbers",
			input:  "[1, 2, 3]",
			want:   []any{1.0, 2.0, 3.0},
			wantOk: true,
		},
		{
			name:   "mixed types",
			input:  `[1, "text", true, null]`,
			want:   []any{1.0, []byte("text"), true, nil},
			wantOk: true,
		},
		{
			name:   "nested arrays",
			input:  "[[1, 2], [3, 4]]",
			want:   []any{[]any{1.0, 2.0}, []any{3.0, 4.0}},
			wantOk: true,
		},
		{
			name:   "with whitespace",
			input:  "[  1  ,  2  ,  3  ]",
			want:   []any{1.0, 2.0, 3.0},
			wantOk: true,
		},
		{
			name:   "trailing comma",
			input:  "[1, 2,]",
			wantOk: false,
		},
		{
			name:   "optional commas",
			input:  "[1 2 3]",
			want:   []any{1.0, 2.0, 3.0},
			wantOk: true,
		},
		{
			name:   "unterminated array",
			input:  "[1, 2",
			wantOk: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.input))
			got, ok := parseAttributeArray(reader)

			if ok != tc.wantOk {
				t.Errorf("parseAttributeArray() ok = %v, want %v", ok, tc.wantOk)
				return
			}

			if !tc.wantOk {
				return
			}

			if len(got) != len(tc.want) {
				t.Errorf("parseAttributeArray() length = %d, want %d", len(got), len(tc.want))
				return
			}

			for i, wantItem := range tc.want {
				compareAttributeValue(t, i, got[i], wantItem)
			}
		})
	}
}

func TestParseAttributeName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   []byte
		wantOk bool
	}{
		{
			name:   "simple name",
			input:  "attr",
			want:   []byte("attr"),
			wantOk: true,
		},
		{
			name:   "name with underscores",
			input:  "my_attr",
			want:   []byte("my_attr"),
			wantOk: true,
		},
		{
			name:   "name with hyphens",
			input:  "my-attr",
			want:   []byte("my-attr"),
			wantOk: true,
		},
		{
			name:   "name with dots",
			input:  "my.attr",
			want:   []byte("my.attr"),
			wantOk: true,
		},
		{
			name:   "name with colons",
			input:  "data:attr",
			want:   []byte("data:attr"),
			wantOk: true,
		},
		{
			name:   "name with numbers",
			input:  "attr123",
			want:   []byte("attr123"),
			wantOk: true,
		},
		{
			name:   "uppercase name",
			input:  "ATTR",
			want:   []byte("ATTR"),
			wantOk: true,
		},
		{
			name:   "mixed case",
			input:  "myAttr",
			want:   []byte("myAttr"),
			wantOk: true,
		},
		{
			name:   "starting with underscore",
			input:  "_attr",
			want:   []byte("_attr"),
			wantOk: true,
		},
		{
			name:   "starting with colon",
			input:  ":attr",
			want:   []byte(":attr"),
			wantOk: true,
		},
		{
			name:   "starting with number",
			input:  "1attr",
			wantOk: false,
		},
		{
			name:   "starting with hyphen",
			input:  "-attr",
			wantOk: false,
		},
		{
			name:   "starting with dot",
			input:  ".attr",
			wantOk: false,
		},
		{
			name:   "empty input",
			input:  "",
			wantOk: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.input))
			got, ok := parseAttributeName(reader)

			if ok != tc.wantOk {
				t.Errorf("parseAttributeName() ok = %v, want %v", ok, tc.wantOk)
				return
			}

			if tc.wantOk && !bytes.Equal(got, tc.want) {
				t.Errorf("parseAttributeName() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseShorthandAttribute(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		shorthandChar byte
		wantName      []byte
		wantValue     []byte
		wantOk        bool
	}{
		{
			name:          "id shorthand",
			input:         "#myid",
			shorthandChar: '#',
			wantName:      []byte("id"),
			wantValue:     []byte("myid"),
			wantOk:        true,
		},
		{
			name:          "class shorthand",
			input:         ".myclass",
			shorthandChar: '.',
			wantName:      []byte("class"),
			wantValue:     []byte("myclass"),
			wantOk:        true,
		},
		{
			name:          "id with hyphens",
			input:         "#my-id",
			shorthandChar: '#',
			wantName:      []byte("id"),
			wantValue:     []byte("my-id"),
			wantOk:        true,
		},
		{
			name:          "id with underscores",
			input:         "#my_id",
			shorthandChar: '#',
			wantName:      []byte("id"),
			wantValue:     []byte("my_id"),
			wantOk:        true,
		},
		{
			name:          "class with dots",
			input:         ".my.class",
			shorthandChar: '.',
			wantName:      []byte("class"),
			wantValue:     []byte("my.class"),
			wantOk:        true,
		},
		{
			name:          "id with colons",
			input:         "#ns:id",
			shorthandChar: '#',
			wantName:      []byte("id"),
			wantValue:     []byte("ns:id"),
			wantOk:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.input))
			attr, ok := parseShorthandAttribute(reader, tc.shorthandChar)

			if ok != tc.wantOk {
				t.Errorf("parseShorthandAttribute() ok = %v, want %v", ok, tc.wantOk)
				return
			}

			if tc.wantOk {
				if !bytes.Equal(attr.Name, tc.wantName) {
					t.Errorf("Name = %q, want %q", attr.Name, tc.wantName)
				}
				if val, ok := attr.Value.([]byte); ok {
					if !bytes.Equal(val, tc.wantValue) {
						t.Errorf("Value = %q, want %q", val, tc.wantValue)
					}
				} else {
					t.Errorf("Value type = %T, want []byte", attr.Value)
				}
			}
		})
	}
}
