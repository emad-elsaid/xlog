package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"xlog"

	"github.com/gin-gonic/gin"
)

func main() {
	cwd, _ := os.Getwd()
	bind := flag.String("bind", "0.0.0.0:7000", "IP and port to bind the web server to")
	source := flag.String("source", cwd, "Directory that will act as a storage")
	flag.Parse()

	absSource, err := filepath.Abs(*source)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Bind: %s, Source: %s", *bind, absSource)
	os.Chdir(absSource)

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	parsedTemplate, err := template.New("").Parse(TEMPLATE)
	if err != nil {
		log.Fatal(err)
	}
	r.SetHTMLTemplate(parsedTemplate)

	r.GET("/:page", pageHandler)
	r.POST("/:page", updateHandler)
	r.GET("/:page/edit", editHandler)
	r.GET("/", pageHandler)
	r.NoRoute(gin.WrapH(http.StripPrefix("/public/", http.FileServer(http.Dir("public")))))

	log.Fatal(r.Run(*bind))
}

func pageHandler(c *gin.Context) {
	page := xlog.NewPage(c.Param("page"))
	if page.Name() == "favicon.ico" {
		return
	}

	if page.Exists() {
		c.HTML(200, "view", gin.H{
			"content": template.HTML(page.Render()),
		})
	} else {
		c.HTML(200, "edit", gin.H{
			"action":  page.Name(),
			"content": page.Content(),
		})
	}
}

func editHandler(c *gin.Context) {
	page := xlog.NewPage(c.Param("page"))

	c.HTML(200, "edit", gin.H{
		"action":  page.Name(),
		"content": page.Content(),
	})
}

func updateHandler(c *gin.Context) {
	page := xlog.NewPage(c.Param("page"))
	content := c.PostForm("content")

	if content != "" {
		page.Write(content)
		c.Redirect(302, "/"+page.Name())
		return
	} else if page.Exists() {
		page.Delete()
	}

	c.Redirect(302, "/")
}
