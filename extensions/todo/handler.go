package todo

import (
	"fmt"
	"regexp"
	"strconv"

	. "github.com/emad-elsaid/xlog"
)

func init() {
	Post(`/\+/todo`, toggleHandler)
}

var taskListRegexp = regexp.MustCompile(`^\[([\sxX])\]\s*`)

func toggleHandler(w Response, r Request) Output {
	if READONLY {
		return Unauthorized("")
	}

	page := NewPage(r.FormValue("page"))
	if !page.Exists() {
		return NotFound(fmt.Sprintf("page: %s not found", page.Name()))
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
	if r.FormValue("checked") == "true" {
		replacement = "[x] "
	}

	line := content[:pos] + taskListRegexp.ReplaceAllString(content[pos:], replacement)
	page.Write(Markdown(line))
	return NoContent()
}
