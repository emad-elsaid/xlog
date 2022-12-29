package rtl

import (
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

const script template.HTML = `
<script>
 // This is a hack to display right to left text correctly
 document.querySelectorAll('.content *')
         .forEach( ele => ele.setAttribute("dir", "auto") );
</script>
`

func init() {
	RegisterWidget(AFTER_VIEW_WIDGET, 1, scriptWidget)
}

func scriptWidget(_ Page) template.HTML {
	return script
}
