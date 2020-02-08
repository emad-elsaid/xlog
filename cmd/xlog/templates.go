package main

const TEMPLATE = `
{{define "view"}}
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8">
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.8.0/css/bulma.css">
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

{{define "edit"}}
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8">
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.8.0/css/bulma.css">
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.css">
		<script src="https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js"></script>
  </head>
	<body>
		<section class="section">
				<form method="POST" action="/{{ .action }}" class="form">
					<div class="field">
						<div class="control has-text-right">
							<button type="submit" class="button is-link">Save</button>
						</div>
					</div>
					<div class="field">
						<textarea name="content" autofocus>{{ .content }}</textarea>
					</div>
					<div class="field">
						<div class="control has-text-right">
							<button type="submit" class="button is-link">Save</button>
						</div>
					</div>
				</form>
		</section>
		<script> var simplemde = new SimpleMDE({autofocus: true, spellChecker: false, toolbar: false}); </script>
	</body>
</html>
{{end}}
`
