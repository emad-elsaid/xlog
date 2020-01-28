package main

const TEMPLATE = `
{{define "style.css"}}
html {
    height: 100%;
    margin: 0;
    padding: 0;
}

body {
    height: 100%;
    margin: 0;
    padding: 0;
    font-family: sans-serif;
}

form.edit {
    height: 100%;
    margin: 0;
    padding: 0;
}

form.edit table {
    height: 100%;
    width: 100%;
    border-collapse: collapse;
    border-spacing: 0;
}

form.edit table td {
    padding: 0;
}


form.edit input[name=title]{
    width: 100%;
    border: 0 none;
    outline: 0 none;
    font-size: 2em;
    padding: 1em;
    font-family: sans-serif;
}

form.edit textarea[name=content] {
    width: 100%;
    height: 100%;
    padding: 1em;
    border: 0 none;
    border-top: 1px solid lightgrey;
    outline: 0 none;
    font-size: 1.2em;
    line-height: 1.8em;
    font-family: sans-serif;
}

.collapse {
    height: 0%;
}

.actions {
    padding: 1em;
    text-align: right;
}

.actions button {
    font-size: 1.2em;
    padding: 0.25em;
}

section.page {
    padding: 1em;
    line-height: 1.6em;
}

a {
    color: inherit;
}
{{end}}
{{define "view.html"}}
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8">
    <title>{{.title}}</title>
    <style>{{template "style.css"}}</style>
  </head>
  <body>
    <section class="page">
      <h1>{{.title}}</h1>
      {{.content}}
    </section>
  </body>
</html>
{{end}}

{{define "edit.html"}}
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8">
    <style>{{template "style.css"}}</style>
  </head>
	<body>
		<form method="POST" action="/{{ .action }}" class="edit">
      <table>
        <tr class="collapse">
          <td>
			      <input name="title" type="text" value="{{ .title }}" autofocus />
          </td>
        </tr>
        <tr>
          <td>
			      <textarea name="content">{{ .content }}</textarea>
          </td>
        </tr>
        <tr class="collapse">
          <td>
            <div class="actions">
              <button type="submit" >Save</button>
            </div>
          </td>
        </tr>
      </table>
		</form>
	</body>
</html>
{{end}}
`
