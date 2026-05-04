package todo

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/emad-elsaid/xlog"
)

const (
	taskUnchecked = "[ ] "
	taskChecked   = "[x] "
)

var taskListRegexp = regexp.MustCompile(`^\[([\sxX])\]\s*`)

func toggleHandler(r xlog.Request) xlog.Output {
	page := xlog.NewPage(r.FormValue("page"))
	if page == nil || !page.Exists() {
		return xlog.NotFound(fmt.Sprintf("page: %s not found", r.FormValue("page")))
	}

	pos, err := strconv.ParseInt(r.FormValue("pos"), 10, 64)
	if err != nil {
		return xlog.BadRequest("Pos value is incorrect, " + err.Error())
	}

	content := string(page.Content())
	if int(pos) >= len(content) {
		return xlog.BadRequest("pos is longer than the content")
	}

	replacement := taskUnchecked
	if len(r.FormValue("checked")) > 0 {
		replacement = taskChecked
	}

	line := content[:pos] + taskListRegexp.ReplaceAllString(content[pos:], replacement)
	page.Write(xlog.Markdown(line))
	return xlog.NoContent()
}
