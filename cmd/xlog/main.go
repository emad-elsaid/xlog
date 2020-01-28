package main

import (
	"flag"
	"html/template"
	"net/http"
	"xlog"

	"github.com/gin-gonic/gin"
)

func main() {
	bind := flag.String("bind", "127.0.0.1:7000", "IP and port to bind the web server to")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	parsedTemplate, err := template.New("").Parse(TEMPLATE)
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(parsedTemplate)

	r.GET("/:page", pageHandler)
	r.POST("/:page", updateHandler)
	r.GET("/:page/edit", editHandler)
	r.GET("/", pageHandler)
	r.NoRoute(gin.WrapH(http.StripPrefix("/public/", http.FileServer(http.Dir("public")))))
	r.Run(*bind)
}

func pageHandler(c *gin.Context) {
	page := xlog.NewPage(c.Param("page"))
	if page.Name() == "favicon.ico" {
		return
	}

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
