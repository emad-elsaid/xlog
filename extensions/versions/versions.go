package versions

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"regexp"

	. "github.com/emad-elsaid/xlog"
)

func init() {
	Listen(BeforeWrite, WriteVersion)
	RegisterProperty(VersionProps)
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
	return os.WriteFile(path.Join(dir, sum), content, 0644)
}

func VersionProps(p Page) []Property {
	dir := p.FileName() + ".versions"
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	i := 0
	for _, f := range files {
		if f.Type().IsRegular() && path.Ext(f.Name()) == ".md" {
			i++
		}
	}

	if i == 0 {
		return nil
	}

	return []Property{prop(i)}
}

type prop int

func (_ prop) Icon() string { return "fa-solid fa-code-branch" }
func (l prop) Name() string { return fmt.Sprintf("%d versions", l) }
