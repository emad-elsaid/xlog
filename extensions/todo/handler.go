package todo

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/emad-elsaid/xlog"
	. "github.com/emad-elsaid/xlog"
)

var taskListRegexp = regexp.MustCompile(`^\[([\sxX])\]\s*`)

func toggleHandler(r Request) Output {
	app := GetApp()
	page := app.NewPage(r.FormValue("page"))
	if page == nil || !page.Exists() {
		return xlog.NotFound(fmt.Sprintf("page: %s not found", r.FormValue("page")))
	}

	pos, err := strconv.Atoi(r.FormValue("pos"))
	if err != nil {
		return xlog.BadRequest("Pos value is incorrect, " + err.Error())
	}

	content := string(page.Content())
	if int(pos) >= len(content) {
		return xlog.BadRequest("pos is longer than the content")
	}

	replacement := "[ ] "
	if len(r.FormValue("checked")) > 0 {
		replacement = "[x] "
	}

	line := content[:pos] + taskListRegexp.ReplaceAllString(content[pos:], replacement)
	page.Write(Markdown(line))
	return xlog.NoContent()
}
