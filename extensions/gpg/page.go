package gpg

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/emad-elsaid/xlog"
	emojiAst "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type page struct {
	name string
	ast  ast.Node
}

func (p *page) Name() string     { return p.name }
func (p *page) Title() string    { return p.name }
func (p *page) FileName() string { return filepath.FromSlash(p.name) + EXT }

func (p *page) GetMeta() (xlog.Metadata, bool) {
	return xlog.Metadata{}, false
}

func (p *page) Exists() bool {
	_, err := os.Stat(p.FileName())
	return err == nil
}

func (p *page) Render() template.HTML {
	content := p.Content()
	content = xlog.PreProcess(content)
	var buf bytes.Buffer
	if err := xlog.MarkDownRenderer.Convert([]byte(content), &buf); err != nil {
		return template.HTML(err.Error())
	}

	return template.HTML(buf.String())
}

func (p *page) Content() xlog.Markdown {
	cmd := exec.Command("gpg", "--decrypt", p.FileName())
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Coudln't decrypt file: %s, err: %s", p.FileName(), err.Error())
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
	defer xlog.Trigger(xlog.AfterDelete, p)

	if p.Exists() {
		err := os.Remove(p.FileName())
		if err != nil {
			fmt.Printf("Can't delete `%s`, err: %s\n", p.Name(), err)
			return false
		}
	}
	return true
}

func (p *page) Write(content xlog.Markdown) bool {
	xlog.Trigger(xlog.BeforeWrite, p)
	defer xlog.Trigger(xlog.AfterWrite, p)

	name := p.FileName()
	os.MkdirAll(filepath.Dir(name), 0700)

	content = xlog.Markdown(strings.ReplaceAll(string(content), "\r\n", "\n"))
	cmd := exec.Command("gpg", "-r", gpgId, "--output", p.FileName(), "--batch", "--yes", "--encrypt")
	cmd.Stdin = bytes.NewBuffer([]byte(content))

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Can't write `%s`, out: %s, err: %s\n", p.Name(), out, err)
		return false
	}

	return true
}

func (p *page) AST() ([]byte, ast.Node) {
	src := p.Content()
	if p.ast == nil {
		p.ast = xlog.MarkDownRenderer.Parser().Parse(text.NewReader([]byte(src)))
	}

	return []byte(src), p.ast
}
func (p *page) Emoji() string {
	_, tree := p.AST()
	if e, ok := xlog.FindInAST[*emojiAst.Emoji](tree); ok {
		return string(e.Value.Unicode)
	}

	return ""
}
