package todo

import (
	"fmt"
	"regexp"
	"strconv"

	. "github.com/emad-elsaid/xlog"
)

var taskListRegexp = regexp.MustCompile(`^\[([\sxX])\]\s*`)

func toggleHandler(r Request) Output {
	page := NewPage(r.FormValue("page"))
	if page == nil || !page.Exists() {
		return NotFound(fmt.Sprintf("page: %s not found", r.FormValue("page")))
	}

	pos, err := strconv.ParseInt(r.FormValue("pos"), 10, 64)
	if err != nil {
		return BadRequest("Pos value is incorrect, " + err.Error())
	}

	content := string(page.Content())
	if int(pos) >= len(content) {
		return BadRequest("pos is longer than the content")
	}

	replacement := "[ ] "
	if len(r.FormValue("checked")) > 0 {
		replacement = "[x] "
	}

	line := content[:pos] + taskListRegexp.ReplaceAllString(content[pos:], replacement)
	page.Write(Markdown(line))
	return NoContent()
}
