package xlog

import (
	"io/fs"

	"github.com/emad-elsaid/types"
)

var extension_page = types.Map[string, bool]{}
var extension_page_enclosed = types.Map[string, bool]{}
var build_perms fs.FileMode = 0744

// RegisterBuildPage registers a path of a page to export when building static version
// of the knowledgebase. encloseInDir will write the output to p/index.html
// instead instead of writing to p directly. that can help have paths with no
// .html extension to be served with the exact name.
func RegisterBuildPage(p string, encloseInDir bool) {
	app := GetApp()
	app.RegisterBuildPage(p, encloseInDir)
}
