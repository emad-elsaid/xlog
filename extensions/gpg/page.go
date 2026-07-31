package gpg

import (
	"bytes"
	"html"
	"html/template"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

type page struct {
	name string
	ast  ast.Node
}

func (p *page) Name() string { return p.name }
func (p *page) FileName() string {
	if !xlog.ValidPageName(p.name) {
		return ""
	}

	return filepath.FromSlash(p.name) + EXT
}

func (p *page) Exists() bool {
	_, err := os.Stat(p.FileName())
	return err == nil
}

func (p *page) Render() template.HTML {
	content := p.Content()
	content = xlog.PreProcess(content)
	var buf bytes.Buffer
	if err := xlog.MarkdownConverter().Convert([]byte(content), &buf); err != nil {
		// #nosec G203 -- Error message is html.EscapeString-sanitized before conversion
		return template.HTML(html.EscapeString(err.Error()))
	}

	// #nosec G203 - buf.String() contains markdown-converted output which is already escaped by MarkdownConverter
	return template.HTML(buf.String())
}

func (p *page) Content() xlog.Markdown {
	cmd := exec.Command("gpg", "--decrypt", p.FileName()) // #nosec G204 -- FileName is from trusted page system, not direct user input
	out, err := cmd.Output()
	if err != nil {
		slog.Error("Coudln't decrypt", "file", p.FileName(), "error", err)
	}

	return xlog.Markdown(out)
}

func (p *page) ModTime() time.Time {
	s, err := os.Stat(p.FileName())
	if err != nil {
		return time.Time{}
	}

	return s.ModTime()
}

func (p *page) Delete() bool {
	defer xlog.Trigger(xlog.PageDeleted, p)

	if p.Exists() {
		err := os.Remove(p.FileName())
		if err != nil {
			slog.Error("Can't delete page", "page", p.Name(), "error", err)
			return false
		}
	}
	return true
}

func (p *page) Write(content xlog.Markdown) bool {
	defer xlog.Trigger(xlog.PageChanged, p)

	name := p.FileName()
	if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
		slog.Error("Can't create directory", "page", p.Name(), "error", err)
		return false
	}

	content = xlog.Markdown(strings.ReplaceAll(string(content), "\r\n", "\n"))
	cmd := exec.Command("gpg", "-r", gpgId, "--output", p.FileName(), "--batch", "--yes", "--encrypt") // #nosec G204 -- FileName is from trusted page system, gpgId from config
	cmd.Stdin = bytes.NewBuffer([]byte(content))

	out, err := cmd.Output()
	if err != nil {
		slog.Error("Can't write page", "page", p.Name(), "output", string(out), "error", err)
		return false
	}

	return true
}

func (p *page) AST() ([]byte, ast.Node) {
	src := p.Content()
	if p.ast == nil {
		p.ast = xlog.MarkdownConverter().Parser().Parse(text.NewReader([]byte(src)))
	}

	return []byte(src), p.ast
}
