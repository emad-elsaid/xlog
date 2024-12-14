package todo

import (
	"fmt"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

const script = `
<script>
(()=>{
	function toggleCheckbox() {
	  const csrf = document.querySelector("input[name=csrf]").value;
	  const data = new FormData();

	  data.append('csrf', csrf);
	  data.append('checked', this.checked);
	  data.append('page', '%s');
	  data.append('pos', this.dataset.pos);

	  fetch("/+/todo", {method: 'POST', body: data});
	}
	let todos = document.querySelectorAll(".view input[type=checkbox][data-pos]");
	todos.forEach(elem => elem.addEventListener("click", toggleCheckbox));
})()
</script>
`

func scriptWidget(p Page) template.HTML {
	return template.HTML(fmt.Sprintf(script, template.JSEscapeString(p.Name())))
}
