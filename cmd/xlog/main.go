package main

import (
	"html/template"
	"net/http"
	"xlog"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/:page", pageHandler)
	r.POST("/:page", updateHandler)

	r.GET("/:page/edit", editHandler)

	r.GET("/", pageHandler)
	r.POST("/", updateHandler)

	r.NoRoute(gin.WrapH(http.StripPrefix("/public/", http.FileServer(http.Dir("public")))))
	r.Run()
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
