package main

const TEMPLATE = `
{{define "style.css"}}
html, body {
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
    width: 100%;
    border-collapse: collapse;
    border-spacing: 0;
}

form.edit table td {
    padding: 0;
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

{{end}}
{{define "view.html"}}
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8">
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.8.0/css/bulma.css">
    <style>{{template "style.css"}}</style>
  </head>
  <body>
    <section class="section">
			<div class="container is-fluid content" dir="auto">
				{{.content}}
			</div>
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
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.css">
		<script src="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js"></script>
  </head>
	<body>
		<form method="POST" action="/{{ .action }}" class="edit">
      <table>
        <tr>
          <td>
			      <textarea name="content" autofocus>{{ .content }}</textarea>
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
		<script> var simplemde = new SimpleMDE({autofocus: true, spellChecker: false}); </script>
	</body>
</html>
{{end}}
`
