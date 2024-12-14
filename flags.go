package xlog

import (
	"flag"
	"os"
)

type Configuration struct {
	Source             string // path to markdown files directory
	Build              string // path to write built files
	Readonly           bool   // is xlog in readonly mode
	Sitename           string // name of knowledgebase
	Index              string // name of the index page markdown file
	NotFoundPage       string // name of the index page markdown file
	BindAddress        string // bind address for the server
	ServeInsecure      bool   // should the server use https for cookie
	CsrfCookieName     string
	DisabledExtensions string
}

var Config Configuration

func init() {
	// Uses current working directory as default value for source flag. If the
	// source flag is set by user the program changes working directory to is
	// and the rest of the program can use relative paths to access files
	cwd, _ := os.Getwd()
	flag.StringVar(&Config.Source, "source", cwd, "Directory that will act as a storage")
	flag.StringVar(&Config.Build, "build", "", "Build all pages as static site in this directory")
	flag.StringVar(&Config.Sitename, "sitename", "XLOG", "Site name is the name that appears on the header beside the logo and in the title tag")
	flag.StringVar(&Config.Index, "index", "index", "Index file name used as home page")
	flag.StringVar(&Config.NotFoundPage, "notfoundpage", "404", "Custom not found page")
	flag.BoolVar(&Config.Readonly, "readonly", false, "Should xlog hide write operations, read-only means all write operations will be disabled")
	flag.StringVar(&Config.BindAddress, "bind", "127.0.0.1:3000", "IP and port to bind the web server to")
	flag.BoolVar(&Config.ServeInsecure, "serve-insecure", false, "Accept http connections and forward crsf cookie over non secure connections")
	flag.StringVar(&Config.CsrfCookieName, "csrf-cookie", "xlog_csrf", "CSRF cookie name")
	flag.StringVar(&Config.DisabledExtensions, "disabled-extensions", "", "disable list of extensions by name, comma separated")
}
