package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/emad-elsaid/xlog"

	_ "embed"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//go:embed template.html
var TEMPLATE string

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

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("xlog-session", store))
	r.Use(csrf.Middleware(csrf.Options{
		Secret: string(xlog.CSRFToken()),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token incorrect")
			c.Abort()
		},
	}))

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
		c.HTML(200, "view", gin.H{
			"content": template.HTML(page.Render()),
		})
	} else {
		c.HTML(200, "edit", gin.H{
			"action":  page.Name(),
			"content": page.Content(),
			"csrf":    csrf.GetToken(c),
		})
	}
}

func editHandler(c *gin.Context) {
	page := xlog.NewPage(c.Param("page"))

	c.HTML(200, "edit", gin.H{
		"action":  page.Name(),
		"content": page.Content(),
		"csrf":    csrf.GetToken(c),
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
