package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func init() {
	PageEvents.Listen(BeforeWrite, WriteVersion)
}

func WriteVersion(p *Page) error {
	if !p.Exists() {
		return nil
	}

	content := []byte(p.Content())
	sum := fmt.Sprintf("%x.md", sha256.Sum256(content))
	dir := p.FileName() + ".versions"

	os.Mkdir(dir, 0700)
	return ioutil.WriteFile(path.Join(dir, sum), content, 0644)
}
