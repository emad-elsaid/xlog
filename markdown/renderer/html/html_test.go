package html

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/text"
)

const htmlEntityLT = "&lt;"

// Test Config and Options

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()
	if cfg.Writer == nil {
		t.Error("NewConfig().Writer should not be nil")
	}
	if cfg.HardWraps != false {
		t.Error("NewConfig().HardWraps should be false")
	}
	if cfg.EastAsianLineBreaks != EastAsianLineBreaksNone {
		t.Error("NewConfig().EastAsianLineBreaks should be EastAsianLineBreaksNone")
	}
	if cfg.XHTML != false {
		t.Error("NewConfig().XHTML should be false")
	}
	if cfg.Unsafe != false {
		t.Error("NewConfig().Unsafe should be false")
	}
}

func TestConfigSetOptionHardWraps(t *testing.T) {
	cfg := &Config{}
	cfg.SetOption(optHardWraps, true)
	if cfg.HardWraps != true {
		t.Error("SetOption(optHardWraps, true) did not set HardWraps")
	}
}

func TestConfigSetOptionEastAsianLineBreaks(t *testing.T) {
	cfg := &Config{}
	cfg.SetOption(optEastAsianLineBreaks, EastAsianLineBreaksSimple)
	if cfg.EastAsianLineBreaks != EastAsianLineBreaksSimple {
		t.Error("SetOption(optEastAsianLineBreaks, EastAsianLineBreaksSimple) did not set value")
	}
}

func TestConfigSetOptionXHTML(t *testing.T) {
	cfg := &Config{}
	cfg.SetOption(optXHTML, true)
	if cfg.XHTML != true {
		t.Error("SetOption(optXHTML, true) did not set XHTML")
	}
}

func TestConfigSetOptionUnsafe(t *testing.T) {
	cfg := &Config{}
	cfg.SetOption(optUnsafe, true)
	if cfg.Unsafe != true {
		t.Error("SetOption(optUnsafe, true) did not set Unsafe")
	}
}

func TestConfigSetOptionWriter(t *testing.T) {
	cfg := &Config{}
	w := NewWriter()
	cfg.SetOption(optTextWriter, w)
	if cfg.Writer != w {
		t.Error("SetOption(optTextWriter, w) did not set Writer")
	}
}

func TestWithWriter(t *testing.T) {
	customWriter := NewWriter()
	opt := WithWriter(customWriter)

	cfg := &Config{}
	opt.SetHTMLOption(cfg)

	if cfg.Writer != customWriter {
		t.Error("WithWriter() did not set custom writer")
	}
}

func TestWithHardWraps(t *testing.T) {
	opt := WithHardWraps()

	cfg := &Config{}
	opt.SetHTMLOption(cfg)
	if cfg.HardWraps != true {
		t.Error("WithHardWraps() did not enable HardWraps")
	}
}

func TestWithEastAsianLineBreaks(t *testing.T) {
	opt := WithEastAsianLineBreaks(EastAsianLineBreaksCSS3Draft)

	cfg := &Config{}
	opt.SetHTMLOption(cfg)
	if cfg.EastAsianLineBreaks != EastAsianLineBreaksCSS3Draft {
		t.Error("WithEastAsianLineBreaks() did not set EastAsianLineBreaks")
	}
}

func TestWithXHTML(t *testing.T) {
	opt := WithXHTML()

	cfg := &Config{}
	opt.SetHTMLOption(cfg)
	if cfg.XHTML != true {
		t.Error("WithXHTML() did not enable XHTML")
	}
}

func TestWithUnsafe(t *testing.T) {
	opt := WithUnsafe()

	cfg := &Config{}
	opt.SetHTMLOption(cfg)
	if cfg.Unsafe != true {
		t.Error("WithUnsafe() did not enable Unsafe")
	}
}

// Test EastAsianLineBreaks

func TestEastAsianLineBreaksNone(t *testing.T) {
	result := EastAsianLineBreaksNone.softLineBreak('好', '的')
	if result != false {
		t.Error("EastAsianLineBreaksNone should always return false")
	}
}

func TestEastAsianLineBreaksSimple(t *testing.T) {
	result := EastAsianLineBreaksSimple.softLineBreak('好', '的')
	if result != false {
		t.Error("EastAsianLineBreaksSimple with both wide should return false")
	}

	result = EastAsianLineBreaksSimple.softLineBreak('a', '的')
	if result != true {
		t.Error("EastAsianLineBreaksSimple with first narrow should return true")
	}
}

func TestEastAsianLineBreaksCSS3DraftZeroWidthSpace(t *testing.T) {
	result := EastAsianLineBreaksCSS3Draft.softLineBreak('\u200B', '的')
	if result != false {
		t.Error("CSS3Draft with zero-width space should return false")
	}
}

func TestEastAsianLineBreaksCSS3DraftWideChars(t *testing.T) {
	result := EastAsianLineBreaksCSS3Draft.softLineBreak('好', '的')
	if result != false {
		t.Error("CSS3Draft with both wide non-Hangul should return false")
	}
}

func TestEastAsianLineBreaksCSS3DraftHangul(t *testing.T) {
	result := EastAsianLineBreaksCSS3Draft.softLineBreak('한', '글')
	if result != true {
		t.Error("CSS3Draft with Hangul should return true")
	}
}

func TestEastAsianLineBreaksCSS3DraftPunctuation(t *testing.T) {
	result := EastAsianLineBreaksCSS3Draft.softLineBreak('.', 'a')
	if result != false {
		t.Error("CSS3Draft with punctuation should return false")
	}
}

func TestEastAsianLineBreaksCSS3DraftNormalChars(t *testing.T) {
	result := EastAsianLineBreaksCSS3Draft.softLineBreak('a', 'b')
	if result != true {
		t.Error("CSS3Draft with normal chars should return true")
	}
}

// Test Renderer

func TestNewRenderer(t *testing.T) {
	r := NewRenderer()
	if r == nil {
		t.Fatal("NewRenderer() returned nil")
	}

	renderer := r.(*Renderer)
	if renderer.Writer == nil {
		t.Error("NewRenderer() Writer is nil")
	}
	if renderer.HardWraps != false {
		t.Error("NewRenderer() HardWraps should be false by default")
	}
}

func TestNewRendererWithOptions(t *testing.T) {
	r := NewRenderer(
		WithHardWraps(),
		WithXHTML(),
		WithUnsafe(),
		WithEastAsianLineBreaks(EastAsianLineBreaksSimple),
	)

	renderer := r.(*Renderer)
	if renderer.HardWraps != true {
		t.Error("WithHardWraps() option not applied")
	}
	if renderer.XHTML != true {
		t.Error("WithXHTML() option not applied")
	}
	if renderer.Unsafe != true {
		t.Error("WithUnsafe() option not applied")
	}
	if renderer.EastAsianLineBreaks != EastAsianLineBreaksSimple {
		t.Error("WithEastAsianLineBreaks() option not applied")
	}
}

func TestRendererRegisterFuncs(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	registered := make(map[ast.NodeKind]bool)
	mockReg := &mockRegisterer{registered: registered}

	renderer.RegisterFuncs(mockReg)

	expectedKinds := []ast.NodeKind{
		ast.KindDocument,
		ast.KindHeading,
		ast.KindBlockquote,
		ast.KindCodeBlock,
		ast.KindFencedCodeBlock,
		ast.KindHTMLBlock,
		ast.KindList,
		ast.KindListItem,
		ast.KindParagraph,
		ast.KindTextBlock,
		ast.KindThematicBreak,
		ast.KindAutoLink,
		ast.KindCodeSpan,
		ast.KindEmphasis,
		ast.KindImage,
		ast.KindLink,
		ast.KindRawHTML,
		ast.KindText,
		ast.KindString,
	}

	for _, kind := range expectedKinds {
		if !registered[kind] {
			t.Errorf("RegisterFuncs() did not register %v", kind)
		}
	}

	if len(registered) != len(expectedKinds) {
		t.Errorf("RegisterFuncs() registered %d kinds, expected %d", len(registered), len(expectedKinds))
	}
}

type mockRegisterer struct {
	registered map[ast.NodeKind]bool
}

func (m *mockRegisterer) Register(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
	m.registered[kind] = true
}

// Test Writer

func TestNewWriter(t *testing.T) {
	w := NewWriter()
	if w == nil {
		t.Fatal("NewWriter() returned nil")
	}
}

func TestNewWriterWithEscapedSpace(t *testing.T) {
	w := NewWriter(WithEscapedSpace())
	if w == nil {
		t.Fatal("NewWriter(WithEscapedSpace()) returned nil")
	}

	dw := w.(*defaultWriter)
	if dw.EscapedSpace != true {
		t.Error("WithEscapedSpace() option not applied")
	}
}

func TestSecureWriteNormal(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter()

	w.SecureWrite(writer, []byte("hello world"))
	_ = writer.Flush()

	if buf.String() != "hello world" {
		t.Errorf("SecureWrite() = %q, want %q", buf.String(), "hello world")
	}
}

func TestSecureWriteNullByte(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter()

	w.SecureWrite(writer, []byte("hello\x00world"))
	_ = writer.Flush()

	if buf.String() != "hello\ufffdworld" {
		t.Errorf("SecureWrite() = %q, want %q", buf.String(), "hello\ufffdworld")
	}
}

func TestRawWriteHTMLChars(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter()

	w.RawWrite(writer, []byte("<>&\""))
	_ = writer.Flush()

	if buf.String() != "&lt;&gt;&amp;&quot;" {
		t.Errorf("RawWrite() = %q, want %q", buf.String(), "&lt;&gt;&amp;&quot;")
	}
}

func TestWriteEscapedPunct(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter()

	w.Write(writer, []byte(`\*\[`))
	_ = writer.Flush()

	if buf.String() != "*[" {
		t.Errorf("Write() = %q, want %q", buf.String(), "*[")
	}
}

func TestWriteNullByte(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter()

	w.Write(writer, []byte("a\x00b"))
	_ = writer.Flush()

	if buf.String() != "a\ufffdb" {
		t.Errorf("Write() = %q, want %q", buf.String(), "a\ufffdb")
	}
}

func TestWriteNumericRefHex(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter()

	w.Write(writer, []byte("&#x3C;"))
	_ = writer.Flush()

	if buf.String() != htmlEntityLT {
		t.Errorf("Write() = %q, want %q", buf.String(), htmlEntityLT)
	}
}

func TestWriteNumericRefDecimal(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter()

	w.Write(writer, []byte("&#60;"))
	_ = writer.Flush()

	if buf.String() != htmlEntityLT {
		t.Errorf("Write() = %q, want %q", buf.String(), htmlEntityLT)
	}
}

func TestWriteEntityReference(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter()

	w.Write(writer, []byte(htmlEntityLT))
	_ = writer.Flush()

	expected := htmlEntityLT
	if buf.String() != expected {
		t.Errorf("Write() = %q, want %q", buf.String(), expected)
	}
}

func TestWriteEscapedSpaceWithoutOption(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter()

	w.Write(writer, []byte(`\ `))
	_ = writer.Flush()

	if buf.String() != `\ ` {
		t.Errorf("Write() = %q, want %q", buf.String(), `\ `)
	}
}

func TestWriteEscapedSpaceWithOption(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	w := NewWriter(WithEscapedSpace())

	w.Write(writer, []byte(`\ `))
	_ = writer.Flush()

	if buf.String() != "" {
		t.Errorf("Write() = %q, want %q", buf.String(), "")
	}
}

// Test IsDangerousURL

func TestIsDangerousURLSafeHTTP(t *testing.T) {
	if IsDangerousURL([]byte("http://example.com")) {
		t.Error("http URL should not be dangerous")
	}
}

func TestIsDangerousURLSafeDataImage(t *testing.T) {
	urls := []string{
		"data:image/png;base64,xxx",
		"data:image/gif;base64,xxx",
		"data:image/jpeg;base64,xxx",
		"data:image/webp;base64,xxx",
		"data:image/svg+xml;base64,xxx",
	}

	for _, url := range urls {
		if IsDangerousURL([]byte(url)) {
			t.Errorf("%s should not be dangerous", url)
		}
	}
}

func TestIsDangerousURLJavaScript(t *testing.T) {
	if !IsDangerousURL([]byte("javascript:alert(1)")) {
		t.Error("javascript: URL should be dangerous")
	}
}

func TestIsDangerousURLVBScript(t *testing.T) {
	if !IsDangerousURL([]byte("vbscript:msgbox(1)")) {
		t.Error("vbscript: URL should be dangerous")
	}
}

func TestIsDangerousURLFile(t *testing.T) {
	if !IsDangerousURL([]byte("file:///etc/passwd")) {
		t.Error("file: URL should be dangerous")
	}
}

func TestIsDangerousURLDataNonImage(t *testing.T) {
	if !IsDangerousURL([]byte("data:text/html,<script>")) {
		t.Error("data:text/html should be dangerous")
	}
}

func TestIsDangerousURLCaseInsensitive(t *testing.T) {
	if !IsDangerousURL([]byte("JavaScript:alert(1)")) {
		t.Error("JavaScript: (mixed case) should be dangerous")
	}

	if IsDangerousURL([]byte("DATA:IMAGE/PNG;base64,xxx")) {
		t.Error("DATA:IMAGE/PNG (mixed case) should not be dangerous")
	}
}

// Test hasPrefix

func TestHasPrefixExactMatch(t *testing.T) {
	if !hasPrefix([]byte("hello"), []byte("hello")) {
		t.Error("exact match should return true")
	}
}

func TestHasPrefixMatch(t *testing.T) {
	if !hasPrefix([]byte("hello world"), []byte("hello")) {
		t.Error("prefix match should return true")
	}
}

func TestHasPrefixCaseInsensitive(t *testing.T) {
	if !hasPrefix([]byte("HELLO"), []byte("hello")) {
		t.Error("case insensitive match should return true")
	}
}

func TestHasPrefixNoMatch(t *testing.T) {
	if hasPrefix([]byte("world"), []byte("hello")) {
		t.Error("no match should return false")
	}
}

// Test render functions

func TestRenderDocumentEntering(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	doc := ast.NewDocument()

	status, err := renderer.renderDocument(writer, []byte("test"), doc, true)
	if err != nil {
		t.Errorf("renderDocument() error = %v", err)
	}
	if status != ast.WalkContinue {
		t.Errorf("renderDocument() status = %v, want WalkContinue", status)
	}

	_ = writer.Flush()
	if buf.Len() != 0 {
		t.Errorf("renderDocument() wrote %d bytes, want 0", buf.Len())
	}
}

func TestRenderHeadingH1Opening(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	heading := ast.NewHeading(1)
	_, _ = renderer.renderHeading(writer, []byte{}, heading, true)
	_ = writer.Flush()

	if buf.String() != "<h1>" {
		t.Errorf("renderHeading() = %q, want %q", buf.String(), "<h1>")
	}
}

func TestRenderHeadingH1Closing(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	heading := ast.NewHeading(1)
	_, _ = renderer.renderHeading(writer, []byte{}, heading, false)
	_ = writer.Flush()

	if buf.String() != "</h1>\n" {
		t.Errorf("renderHeading() = %q, want %q", buf.String(), "</h1>\n")
	}
}

func TestRenderBlockquoteOpening(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	blockquote := ast.NewBlockquote()
	_, _ = renderer.renderBlockquote(writer, []byte{}, blockquote, true)
	_ = writer.Flush()

	if buf.String() != "<blockquote>\n" {
		t.Errorf("renderBlockquote() = %q, want %q", buf.String(), "<blockquote>\n")
	}
}

func TestRenderBlockquoteClosing(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	blockquote := ast.NewBlockquote()
	_, _ = renderer.renderBlockquote(writer, []byte{}, blockquote, false)
	_ = writer.Flush()

	if buf.String() != "</blockquote>\n" {
		t.Errorf("renderBlockquote() = %q, want %q", buf.String(), "</blockquote>\n")
	}
}

func TestRenderCodeBlock(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	codeBlock := ast.NewCodeBlock()
	codeBlock.Lines().Append(text.NewSegment(0, 4))
	source := []byte("code")

	_, _ = renderer.renderCodeBlock(writer, source, codeBlock, true)
	_, _ = renderer.renderCodeBlock(writer, source, codeBlock, false)
	_ = writer.Flush()

	expected := "<pre><code>code</code></pre>\n"
	if buf.String() != expected {
		t.Errorf("renderCodeBlock() = %q, want %q", buf.String(), expected)
	}
}

func TestRenderThematicBreakHTML5(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	hr := ast.NewThematicBreak()
	_, _ = renderer.renderThematicBreak(writer, []byte{}, hr, true)
	_ = writer.Flush()

	if buf.String() != "<hr>\n" {
		t.Errorf("renderThematicBreak() = %q, want %q", buf.String(), "<hr>\n")
	}
}

func TestRenderThematicBreakXHTML(t *testing.T) {
	r := NewRenderer(WithXHTML())
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	hr := ast.NewThematicBreak()
	_, _ = renderer.renderThematicBreak(writer, []byte{}, hr, true)
	_ = writer.Flush()

	if buf.String() != "<hr />\n" {
		t.Errorf("renderThematicBreak(XHTML) = %q, want %q", buf.String(), "<hr />\n")
	}
}

func TestRenderEmphasisItalic(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	emphasis := ast.NewEmphasis(1)
	_, _ = renderer.renderEmphasis(writer, []byte{}, emphasis, true)
	_ = writer.Flush()

	if buf.String() != "<em>" {
		t.Errorf("renderEmphasis(1) = %q, want %q", buf.String(), "<em>")
	}
}

func TestRenderEmphasisBold(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	emphasis := ast.NewEmphasis(2)
	_, _ = renderer.renderEmphasis(writer, []byte{}, emphasis, true)
	_ = writer.Flush()

	if buf.String() != "<strong>" {
		t.Errorf("renderEmphasis(2) = %q, want %q", buf.String(), "<strong>")
	}
}

func TestRenderStringNormal(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	str := ast.NewString([]byte("hello"))
	_, _ = renderer.renderString(writer, []byte{}, str, true)
	_ = writer.Flush()

	if buf.String() != "hello" {
		t.Errorf("renderString() = %q, want %q", buf.String(), "hello")
	}
}

func TestRenderStringCode(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	str := ast.NewString([]byte("<code>"))
	str.SetCode(true)
	_, _ = renderer.renderString(writer, []byte{}, str, true)
	_ = writer.Flush()

	if buf.String() != "<code>" {
		t.Errorf("renderString(code) = %q, want %q", buf.String(), "<code>")
	}
}

func TestRenderStringRaw(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	str := ast.NewString([]byte("<tag>"))
	str.SetRaw(true)
	_, _ = renderer.renderString(writer, []byte{}, str, true)
	_ = writer.Flush()

	if buf.String() != "&lt;tag&gt;" {
		t.Errorf("renderString(raw) = %q, want %q", buf.String(), "&lt;tag&gt;")
	}
}

// TestRenderParagraph tests paragraph rendering.
func TestRenderParagraph(t *testing.T) {
	tests := []struct {
		name     string
		attrs    bool
		entering bool
		want     string
	}{
		{"entering without attrs", false, true, "<p>"},
		{"exiting", false, false, "</p>\n"},
		{"entering with attrs", true, true, "<p"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRenderer()
			renderer := r.(*Renderer)
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			p := ast.NewParagraph()
			if tc.attrs {
				p.SetAttribute([]byte("class"), []byte("test"))
			}

			_, _ = renderer.renderParagraph(writer, []byte{}, p, tc.entering)
			_ = writer.Flush()

			result := buf.String()
			if tc.attrs && tc.entering {
				if !bytes.HasPrefix([]byte(result), []byte(tc.want)) {
					t.Errorf("renderParagraph() = %q, want prefix %q", result, tc.want)
				}
			} else if result != tc.want {
				t.Errorf("renderParagraph() = %q, want %q", result, tc.want)
			}
		})
	}
}

// TestRenderTextBlock tests text block rendering.
func TestRenderTextBlock(t *testing.T) {
	tests := []struct {
		name       string
		entering   bool
		hasNextSib bool
		hasChild   bool
		want       string
	}{
		{"entering", true, false, false, ""},
		{"exiting no sibling", false, false, false, ""},
		{"exiting with sibling and child", false, true, true, "\n"},
		{"exiting with sibling no child", false, true, false, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRenderer()
			renderer := r.(*Renderer)
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			tb := ast.NewTextBlock()
			if tc.hasChild {
				tb.AppendChild(tb, ast.NewText())
			}
			if tc.hasNextSib {
				parent := ast.NewParagraph()
				parent.AppendChild(parent, tb)
				parent.AppendChild(parent, ast.NewText())
			}

			_, _ = renderer.renderTextBlock(writer, []byte{}, tb, tc.entering)
			_ = writer.Flush()

			if buf.String() != tc.want {
				t.Errorf("renderTextBlock() = %q, want %q", buf.String(), tc.want)
			}
		})
	}
}

// TestRenderAutoLink tests auto link rendering.
func TestRenderAutoLink(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		linkType    ast.AutoLinkType
		hasProtocol bool
		wantContain string
	}{
		{
			name:        "regular URL",
			url:         "https://example.com",
			linkType:    ast.AutoLinkURL,
			wantContain: `<a href="https://example.com">`,
		},
		{
			name:        "email without mailto",
			url:         "user@example.com",
			linkType:    ast.AutoLinkEmail,
			wantContain: `<a href="mailto:user@example.com">`,
		},
		{
			name:        "email with mailto protocol",
			url:         "user@example.com",
			linkType:    ast.AutoLinkEmail,
			hasProtocol: true,
			wantContain: `href="mailto:`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRenderer()
			renderer := r.(*Renderer)
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			source := []byte(tc.url)
			txt := ast.NewText()
			txt.Segment = text.NewSegment(0, len(tc.url))

			al := ast.NewAutoLink(tc.linkType, txt)
			if tc.hasProtocol {
				al.Protocol = []byte("mailto")
			}

			_, _ = renderer.renderAutoLink(writer, source, al, true)
			_ = writer.Flush()

			if !bytes.Contains(buf.Bytes(), []byte(tc.wantContain)) {
				t.Errorf("renderAutoLink() = %q, want to contain %q", buf.String(), tc.wantContain)
			}
		})
	}
}

// TestRenderCodeSpan tests code span rendering.
func TestRenderCodeSpan(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		entering bool
		want     string
	}{
		{"entering", "code", true, "<code>code"},
		{"exiting", "", false, "</code>"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRenderer()
			renderer := r.(*Renderer)
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			cs := ast.NewCodeSpan()
			if tc.text != "" {
				source := []byte(tc.text)
				txt := ast.NewTextSegment(text.NewSegment(0, len(source)))
				cs.AppendChild(cs, txt)
				_, _ = renderer.renderCodeSpan(writer, source, cs, tc.entering)
			} else {
				_, _ = renderer.renderCodeSpan(writer, []byte{}, cs, tc.entering)
			}
			_ = writer.Flush()

			if buf.String() != tc.want {
				t.Errorf("renderCodeSpan() = %q, want %q", buf.String(), tc.want)
			}
		})
	}
}

// TestRenderLink tests link rendering.
func TestRenderLink(t *testing.T) {
	tests := []struct {
		name     string
		dest     string
		title    string
		entering bool
		unsafe   bool
		wantPart string
	}{
		{
			name:     "entering safe URL",
			dest:     "https://example.com",
			title:    "Example",
			entering: true,
			unsafe:   false,
			wantPart: `<a href="https://example.com" title="Example">`,
		},
		{
			name:     "exiting",
			entering: false,
			wantPart: `</a>`,
		},
		{
			name:     "dangerous URL blocked",
			dest:     "javascript:alert(1)",
			entering: true,
			unsafe:   false,
			wantPart: `<a href="">`,
		},
		{
			name:     "dangerous URL allowed when unsafe",
			dest:     "javascript:alert(1)",
			entering: true,
			unsafe:   true,
			wantPart: `<a href="javascript:alert(1)">`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRenderer(WithUnsafe())
			renderer := r.(*Renderer)
			renderer.Unsafe = tc.unsafe

			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			link := ast.NewLink()
			link.Destination = []byte(tc.dest)
			if tc.title != "" {
				link.Title = []byte(tc.title)
			}

			_, _ = renderer.renderLink(writer, []byte{}, link, tc.entering)
			_ = writer.Flush()

			if !bytes.Contains(buf.Bytes(), []byte(tc.wantPart)) {
				t.Errorf("renderLink() = %q, want to contain %q", buf.String(), tc.wantPart)
			}
		})
	}
}

// TestRenderImage tests image rendering.
func TestRenderImage(t *testing.T) {
	tests := []struct {
		name     string
		dest     string
		title    string
		alt      string
		xhtml    bool
		unsafe   bool
		wantPart string
	}{
		{
			name:     "basic image",
			dest:     "/img/test.png",
			title:    "Test Image",
			alt:      "test",
			xhtml:    false,
			wantPart: `<img src="/img/test.png" alt="test" title="Test Image">`,
		},
		{
			name:     "XHTML image",
			dest:     "/img/test.png",
			xhtml:    true,
			wantPart: `<img src="/img/test.png" alt="" />`,
		},
		{
			name:     "dangerous URL blocked",
			dest:     "javascript:alert(1)",
			unsafe:   false,
			wantPart: `<img src="" alt="">`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var opts []Option
			if tc.xhtml {
				opts = append(opts, WithXHTML())
			}
			if tc.unsafe {
				opts = append(opts, WithUnsafe())
			}
			r := NewRenderer(opts...)
			renderer := r.(*Renderer)

			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			link := ast.NewLink()
			link.Destination = []byte(tc.dest)
			if tc.title != "" {
				link.Title = []byte(tc.title)
			}

			img := ast.NewImage(link)
			if tc.alt != "" {
				txt := ast.NewText()
				txt.Segment = text.NewSegment(0, len(tc.alt))
				img.AppendChild(img, txt)
			}

			source := []byte(tc.alt)
			_, _ = renderer.renderImage(writer, source, img, true)
			_ = writer.Flush()

			if !bytes.Contains(buf.Bytes(), []byte(tc.wantPart)) {
				t.Errorf("renderImage() = %q, want to contain %q", buf.String(), tc.wantPart)
			}
		})
	}
}

// TestRenderRawHTML tests raw HTML rendering.
func TestRenderRawHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		unsafe   bool
		entering bool
		want     string
	}{
		{
			name:     "unsafe mode entering",
			html:     "<div>test</div>",
			unsafe:   true,
			entering: true,
			want:     "<div>test</div>",
		},
		{
			name:     "safe mode entering",
			html:     "<div>test</div>",
			unsafe:   false,
			entering: true,
			want:     "<!-- raw HTML omitted -->",
		},
		{
			name:     "exiting",
			html:     "",
			unsafe:   false,
			entering: false,
			want:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var opts []Option
			if tc.unsafe {
				opts = append(opts, WithUnsafe())
			}
			r := NewRenderer(opts...)
			renderer := r.(*Renderer)

			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			source := []byte(tc.html)
			raw := ast.NewRawHTML()
			if tc.html != "" {
				raw.Segments.Append(text.NewSegment(0, len(source)))
			}

			_, _ = renderer.renderRawHTML(writer, source, raw, tc.entering)
			_ = writer.Flush()

			if buf.String() != tc.want {
				t.Errorf("renderRawHTML() = %q, want %q", buf.String(), tc.want)
			}
		})
	}
}

// TestRenderText tests text rendering.
func TestRenderText(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		raw          bool
		hardBreak    bool
		softBreak    bool
		hardWraps    bool
		xhtml        bool
		entering     bool
		wantContains string
	}{
		{
			name:         "normal text entering",
			text:         "hello world",
			entering:     true,
			wantContains: "hello world",
		},
		{
			name:     "exiting",
			entering: false,
		},
		{
			name:         "raw text",
			text:         "plain",
			raw:          true,
			entering:     true,
			wantContains: "plain",
		},
		{
			name:         "hard line break",
			text:         "text",
			hardBreak:    true,
			entering:     true,
			wantContains: "<br>",
		},
		{
			name:         "hard line break XHTML",
			text:         "text",
			hardBreak:    true,
			xhtml:        true,
			entering:     true,
			wantContains: "<br />",
		},
		{
			name:         "soft break with hardwraps",
			text:         "text",
			softBreak:    true,
			hardWraps:    true,
			entering:     true,
			wantContains: "<br>",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var opts []Option
			if tc.xhtml {
				opts = append(opts, WithXHTML())
			}
			if tc.hardWraps {
				opts = append(opts, WithHardWraps())
			}
			r := NewRenderer(opts...)
			renderer := r.(*Renderer)

			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			source := []byte(tc.text)
			txt := ast.NewText()
			if tc.text != "" {
				txt.Segment = text.NewSegment(0, len(source))
			}
			if tc.raw {
				txt.SetRaw(true)
			}
			if tc.hardBreak {
				txt.SetHardLineBreak(true)
			}
			if tc.softBreak {
				txt.SetSoftLineBreak(true)
			}

			_, _ = renderer.renderText(writer, source, txt, tc.entering)
			_ = writer.Flush()

			if tc.entering && tc.wantContains != "" {
				if !bytes.Contains(buf.Bytes(), []byte(tc.wantContains)) {
					t.Errorf("renderText() = %q, want to contain %q", buf.String(), tc.wantContains)
				}
			}
		})
	}
}

// TestRenderTexts tests renderTexts helper.
func TestRenderTexts(t *testing.T) {
	r := NewRenderer()
	renderer := r.(*Renderer)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	parent := ast.NewParagraph()
	source := []byte("hello world")

	txt := ast.NewText()
	txt.Segment = text.NewSegment(0, 5)
	parent.AppendChild(parent, txt)

	str := ast.NewString([]byte(" world"))
	parent.AppendChild(parent, str)

	renderer.renderTexts(writer, source, parent)
	_ = writer.Flush()

	if !bytes.Contains(buf.Bytes(), []byte("hello")) {
		t.Errorf("renderTexts() = %q, want to contain 'hello'", buf.String())
	}
}

// TestRenderList tests list rendering.
func TestRenderList(t *testing.T) {
	tests := []struct {
		name     string
		ordered  bool
		start    int
		entering bool
		want     string
	}{
		{
			name:     "unordered entering",
			ordered:  false,
			entering: true,
			want:     "<ul>\n",
		},
		{
			name:     "unordered exiting",
			ordered:  false,
			entering: false,
			want:     "</ul>\n",
		},
		{
			name:     "ordered entering default start",
			ordered:  true,
			start:    1,
			entering: true,
			want:     "<ol>\n",
		},
		{
			name:     "ordered entering custom start",
			ordered:  true,
			start:    5,
			entering: true,
			want:     `<ol start="5"`,
		},
		{
			name:     "ordered exiting",
			ordered:  true,
			entering: false,
			want:     "</ol>\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRenderer()
			renderer := r.(*Renderer)

			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			list := ast.NewList(0)
			if tc.ordered {
				list.Marker = '.'
			}
			list.Start = tc.start

			_, _ = renderer.renderList(writer, []byte{}, list, tc.entering)
			_ = writer.Flush()

			result := buf.String()
			if !bytes.Contains([]byte(result), []byte(tc.want)) {
				t.Errorf("renderList() = %q, want to contain %q", result, tc.want)
			}
		})
	}
}

// TestRenderListItem tests list item rendering.
func TestRenderListItem(t *testing.T) {
	tests := []struct {
		name       string
		hasChild   bool
		textBlock  bool
		entering   bool
		wantPrefix string
	}{
		{
			name:       "entering without attrs",
			entering:   true,
			wantPrefix: "<li>",
		},
		{
			name:     "exiting",
			entering: false,
		},
		{
			name:       "with text block child",
			hasChild:   true,
			textBlock:  true,
			entering:   true,
			wantPrefix: "<li>",
		},
		{
			name:       "with paragraph child",
			hasChild:   true,
			textBlock:  false,
			entering:   true,
			wantPrefix: "<li>\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRenderer()
			renderer := r.(*Renderer)

			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			li := ast.NewListItem(0)
			if tc.hasChild {
				if tc.textBlock {
					li.AppendChild(li, ast.NewTextBlock())
				} else {
					li.AppendChild(li, ast.NewParagraph())
				}
			}

			_, _ = renderer.renderListItem(writer, []byte{}, li, tc.entering)
			_ = writer.Flush()

			result := buf.String()
			if tc.entering && tc.wantPrefix != "" {
				if !bytes.HasPrefix([]byte(result), []byte(tc.wantPrefix)) {
					t.Errorf("renderListItem() = %q, want prefix %q", result, tc.wantPrefix)
				}
			} else if !tc.entering {
				if result != "</li>\n" {
					t.Errorf("renderListItem() exiting = %q, want %q", result, "</li>\n")
				}
			}
		})
	}
}

// TestRenderHTMLBlock tests HTML block rendering.
func TestRenderHTMLBlock(t *testing.T) {
	tests := []struct {
		name       string
		html       string
		hasClosure bool
		unsafe     bool
		entering   bool
		want       string
	}{
		{
			name:     "unsafe entering",
			html:     "<div>\n",
			unsafe:   true,
			entering: true,
			want:     "<div>\n",
		},
		{
			name:     "safe entering",
			html:     "<div>\n",
			unsafe:   false,
			entering: true,
			want:     "<!-- raw HTML omitted -->\n",
		},
		{
			name:       "unsafe exiting with closure",
			html:       "</div>\n",
			hasClosure: true,
			unsafe:     true,
			entering:   false,
			want:       "</div>\n",
		},
		{
			name:       "safe exiting with closure",
			html:       "</div>\n",
			hasClosure: true,
			unsafe:     false,
			entering:   false,
			want:       "<!-- raw HTML omitted -->\n",
		},
		{
			name:     "exiting without closure",
			entering: false,
			want:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var opts []Option
			if tc.unsafe {
				opts = append(opts, WithUnsafe())
			}
			r := NewRenderer(opts...)
			renderer := r.(*Renderer)

			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)

			source := []byte(tc.html)
			block := ast.NewHTMLBlock(ast.HTMLBlockType1)
			if tc.entering && tc.html != "" {
				block.Lines().Append(text.NewSegment(0, len(source)))
			}
			if tc.hasClosure {
				block.ClosureLine = text.NewSegment(0, len(source))
			}

			_, _ = renderer.renderHTMLBlock(writer, source, block, tc.entering)
			_ = writer.Flush()

			if buf.String() != tc.want {
				t.Errorf("renderHTMLBlock() = %q, want %q", buf.String(), tc.want)
			}
		})
	}
}
