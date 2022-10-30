package versions

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	. "github.com/emad-elsaid/xlog"
)

func init() {
	Listen(BeforeWrite, WriteVersion)
	RegisterWidget(ACTION_WIDGET, VersionMeta)
	IgnoreDirectory(regexp.MustCompile(`\.versions$`))
}

func WriteVersion(p Page) error {
	if !p.Exists() {
		return nil
	}

	content := []byte(p.Content())
	sum := fmt.Sprintf("%x.md", sha256.Sum256(content))
	dir := p.FileName() + ".versions"

	os.Mkdir(dir, 0700)
	return ioutil.WriteFile(path.Join(dir, sum), content, 0644)
}

func VersionMeta(p Page, _ Request) (t template.HTML) {
	dir := p.FileName() + ".versions"
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	i := 0
	for _, f := range files {
		if f.Type().IsRegular() && path.Ext(f.Name()) == ".md" {
			i++
		}
	}

	if i == 0 {
		return
	}

	return template.HTML(
		fmt.Sprintf(
			`<span class="icon"><i class="fa-solid fa-code-branch"></i></span><span>%d versions</span>`,
			i,
		),
	)
}
