package autolink_pages

import (
	"context"
	"embed"
	"html/template"
	"path"
	"sort"
	"strings"
	"sync"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	east "github.com/emad-elsaid/xlog/markdown/extension/ast"
)

//go:embed templates
var templates embed.FS

type NormalizedPage struct {
	page           Page
	normalizedName string
}

type fileInfoByNameLength []*NormalizedPage

func (a fileInfoByNameLength) Len() int      { return len(a) }
func (a fileInfoByNameLength) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a fileInfoByNameLength) Less(i, j int) bool {
	return len(a[i].normalizedName) > len(a[j].normalizedName)
}

var autolinkPages []*NormalizedPage
var autolinkPage_lck sync.Mutex

func UpdatePagesList(Page) (err error) {
	autolinkPage_lck.Lock()
	defer autolinkPage_lck.Unlock()

	ps := MapPage(context.Background(), func(p Page) *NormalizedPage {
		return &NormalizedPage{
			page:           p,
			normalizedName: path.Base(strings.ToLower(p.Name())),
		}
	})
	sort.Sort(fileInfoByNameLength(ps))
	autolinkPages = ps
	return
}

func countTodos(p Page) (total int, done int) {
	_, tree := p.AST()
	tasks := FindAllInAST[*east.TaskCheckBox](tree)
	for _, v := range tasks {
		total++
		if v.IsChecked {
			done++
		}
	}

	return
}

func backlinksSection(p Page) template.HTML {
	if p.Name() == Config.Index {
		return ""
	}

	pages := MapPage(context.Background(), func(a Page) Page {
		_, tree := a.AST()
		if a.Name() == p.Name() || !containLinkToFrom(tree, a, p) {
			return nil
		}

		return a
	})

	return Partial("backlinks", Locals{"pages": pages})
}

// containLinkTo checks if an AST node contains a link to the given page.
// For relative links, it matches on basename only, which means links may
// incorrectly match pages with the same basename in different folders.
// This is a known limitation when the source page context is unavailable.
func containLinkTo(n ast.Node, p Page) bool {
	if n.Kind() == KindPageLink {
		t, _ := n.(*PageLink)
		if t.page.FileName() == p.FileName() {
			return true
		}
	}
	if n.Kind() == ast.KindLink {
		t, _ := n.(*ast.Link)
		dst := string(t.Destination)

		// link is absolute: remove / and match full path
		if strings.HasPrefix(dst, "/") {
			path := strings.TrimPrefix(dst, "/")
			if string(path) == p.Name() {
				return true
			}
		} else { // link is relative: match on basename
			// Note: This may incorrectly match pages with the same basename
			// in different folders. Without knowing the source page location,
			// we cannot resolve relative paths accurately.
			base := path.Base(p.Name())
			if dst == base {
				return true
			}
		}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if containLinkTo(c, p) {
			return true
		}

		if c == n.LastChild() {
			break
		}
	}

	return false
}

// containLinkToFrom checks if an AST node contains a link from sourcePage to targetPage.
// This version is aware of the source page context and can better resolve relative links.
func containLinkToFrom(n ast.Node, sourcePage, targetPage Page) bool {
	if n.Kind() == KindPageLink {
		t, _ := n.(*PageLink)
		if t.page.FileName() == targetPage.FileName() {
			return true
		}
	}
	if n.Kind() == ast.KindLink {
		t, _ := n.(*ast.Link)
		dst := string(t.Destination)

		// link is absolute: remove / and match full path
		if strings.HasPrefix(dst, "/") {
			cleanPath := strings.TrimPrefix(dst, "/")
			if cleanPath == targetPage.Name() {
				return true
			}
		} else { // link is relative: resolve from source page's directory
			sourceDir := path.Dir(sourcePage.Name())
			resolvedPath := path.Join(sourceDir, dst)
			// Normalize path by cleaning it
			resolvedPath = path.Clean(resolvedPath)
			
			if resolvedPath == targetPage.Name() {
				return true
			}
			
			// Fallback: also check basename for compatibility
			// This handles cases where relative links use just the filename
			if dst == path.Base(targetPage.Name()) {
				return true
			}
		}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if containLinkToFrom(c, sourcePage, targetPage) {
			return true
		}

		if c == n.LastChild() {
			break
		}
	}

	return false
}
