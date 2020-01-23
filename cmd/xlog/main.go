package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/:page", pageHandler)
	r.POST("/:page", updateHandler)

	r.GET("/:page/e", editHandler)

	r.GET("/", pageHandler)
	r.POST("/", updateHandler)

	r.NoRoute(gin.WrapH(http.StripPrefix("/public/", http.FileServer(http.Dir("public")))))
	r.Run()
}

func pageHandler(c *gin.Context) {
	page := normalizePage(c.Param("page"))

	if pageExists(page) {
		c.HTML(200, "view.html", gin.H{
			"title":   pageTitle(page),
			"content": template.HTML(renderPage(page)),
		})
	} else {
		c.HTML(200, "edit.html", gin.H{
			"action":  page,
			"title":   pageTitle(page),
			"content": pageContent(page),
		})
	}
}

func editHandler(c *gin.Context) {
	page := normalizePage(c.Param("page"))

	c.HTML(200, "edit.html", gin.H{
		"action":  page,
		"title":   pageTitle(page),
		"content": pageContent(page),
	})
}

func updateHandler(c *gin.Context) {
	page := normalizePage(c.Param("page"))
	title := c.PostForm("title")
	content := c.PostForm("content")

	writePage(page, title, content)
	c.Redirect(302, "/"+page)
}

func pageExists(p string) bool {
	filename := p + ".md"
	_, err := os.Stat(filename)
	return err == nil
}

func renderPage(p string) string {
	content := pageContent(p)

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			highlighting.Highlighting,
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return err.Error()
	}

	return buf.String()
}

func pageTitle(p string) string {
	dat, err := ioutil.ReadFile(p + ".md")
	if err != nil {
		fmt.Printf("Can't open `%s`, err: %s\n", p, err)
		return ""
	}

	fileContent := string(dat)
	endOfLine := strings.Index(fileContent, "\n")
	title := fileContent[:endOfLine]
	return title
}

func pageContent(p string) string {
	dat, err := ioutil.ReadFile(p + ".md")
	if err != nil {
		fmt.Printf("Can't open `%s`, err: %s\n", p, err)
	}

	fileContent := string(dat)
	endOfLine := strings.Index(fileContent, "\n")

	fileContent = fileContent[endOfLine+1:]
	endOfLine = strings.Index(fileContent, "\n")

	content := fileContent[endOfLine+1:]
	return content
}

func writePage(page, title, content string) {
	err := ioutil.WriteFile(page+".md", []byte(title+"\n=========\n"+content), 0644)
	if err != nil {
		panic(err)
	}
}

func normalizePage(page string) string {
	if page == "" {
		return "index"
	}

	return page
}
