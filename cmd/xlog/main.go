package main

import (
	"flag"
	"html/template"
	"net/http"
	"xlog"

	"github.com/gin-gonic/gin"
)

const (
	viewTemplate = `
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8">
    <link href="/public/style.css" rel="stylesheet">
    <title>{{.title}}</title>
  </head>
  <body>
    <section class="page">
      <h1>{{.title}}</h1>
      {{.content}}
    </section>
  </body>
</html>
`
	editTemplate = `
<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta http-equiv="Content-Type" content="text/html;charset=utf-8">
    <link href="/public/style.css" rel="stylesheet">
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
`
)

func main() {
	bind := flag.String("bind", "127.0.0.1:7000", "IP and port to bind the web server to")
	flag.Parse()

	r := gin.Default()

	view, err := template.New("view.html").Parse(viewTemplate)
	if err != nil {
		panic(err)
	}
	edit, err := view.New("edit.html").Parse(editTemplate)
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(view)
	r.SetHTMLTemplate(edit)

	r.GET("/:page", pageHandler)
	r.POST("/:page", updateHandler)

	r.GET("/:page/edit", editHandler)

	r.GET("/", pageHandler)

	r.NoRoute(gin.WrapH(http.StripPrefix("/public/", http.FileServer(http.Dir("public")))))
	r.Run(*bind)
}

func pageHandler(c *gin.Context) {
	page := xlog.NewPage(c.Param("page"))

	if page.Exists() {
		c.HTML(200, "view.html", gin.H{
			"title":   page.Title(),
			"content": template.HTML(page.Render()),
		})
	} else {
		c.HTML(200, "edit.html", gin.H{
			"action":  page.Name(),
			"title":   page.Title(),
			"content": page.Content(),
		})
	}
}

func editHandler(c *gin.Context) {
	page := xlog.NewPage(c.Param("page"))

	c.HTML(200, "edit.html", gin.H{
		"action":  page.Name(),
		"title":   page.Title(),
		"content": page.Content(),
	})
}

func updateHandler(c *gin.Context) {
	page := xlog.NewPage(c.Param("page"))
	title := c.PostForm("title")
	content := c.PostForm("content")

	if content != "" {
		page.Write(title, content)
		c.Redirect(302, "/"+page.Name())
		return
	} else if page.Exists() {
		page.Delete()
	}

	c.Redirect(302, "/")
}
