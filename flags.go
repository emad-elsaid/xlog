package xlog

import (
	"flag"
	"os"
)

var (
	SOURCE   string
	BUILD    string
	READONLY bool
	SITENAME string
	SIDEBAR  bool
	INDEX    string
	TEMPLATE string
)

func init() {
	// Uses current working directory as default value for source flag. If the
	// source flag is set by user the program changes working directory to is
	// and the rest of the program can use relative paths to access files
	cwd, _ := os.Getwd()
	flag.StringVar(&SOURCE, "source", cwd, "Directory that will act as a storage")
	flag.StringVar(&BUILD, "build", "", "Build all pages as static site in this directory")
	flag.StringVar(&SITENAME, "sitename", "XLOG", "Site name is the name that appears on the header beside the logo and in the title tag")
	flag.StringVar(&INDEX, "index", "index", "Index file name used as home page")
	flag.StringVar(&TEMPLATE, "template", "template", "Template file name, to use as default content when creating a new page")
	flag.BoolVar(&READONLY, "readonly", false, "Should xlog hide write operations, read-only means all write operations will be disabled")
	flag.BoolVar(&SIDEBAR, "sidebar", true, "Should render sidebar.")
}
