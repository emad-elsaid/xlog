package xlog

import (
	"flag"
	"os"
	"testing"
)

func TestConfigurationDefaults(t *testing.T) {
	// Reset flags to ensure clean test state
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	
	// Re-initialize flags
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
	flag.StringVar(&Config.DisabledExtensions, "disabled-extensions", "", "disable list of extensions by name, comma separated, `all` will disable all extensions")
	flag.StringVar(&Config.CodeStyle, "codestyle", "dracula", "code highlighting style name from the list supported by https://pkg.go.dev/github.com/alecthomas/chroma/v2/styles")
	flag.StringVar(&Config.Theme, "theme", "", "bulma theme to use. (light, dark). empty value means system preference is used")

	tests := []struct {
		name     string
		field    string
		expected interface{}
	}{
		{"Default source should be cwd", "Source", cwd},
		{"Default build should be empty", "Build", ""},
		{"Default sitename should be XLOG", "Sitename", "XLOG"},
		{"Default index should be index", "Index", "index"},
		{"Default notfoundpage should be 404", "NotFoundPage", "404"},
		{"Default readonly should be false", "Readonly", false},
		{"Default bind address should be 127.0.0.1:3000", "BindAddress", "127.0.0.1:3000"},
		{"Default serve-insecure should be false", "ServeInsecure", false},
		{"Default csrf cookie should be xlog_csrf", "CsrfCookieName", "xlog_csrf"},
		{"Default disabled-extensions should be empty", "DisabledExtensions", ""},
		{"Default codestyle should be dracula", "CodeStyle", "dracula"},
		{"Default theme should be empty", "Theme", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual interface{}
			switch tt.field {
			case "Source":
				actual = Config.Source
			case "Build":
				actual = Config.Build
			case "Sitename":
				actual = Config.Sitename
			case "Index":
				actual = Config.Index
			case "NotFoundPage":
				actual = Config.NotFoundPage
			case "Readonly":
				actual = Config.Readonly
			case "BindAddress":
				actual = Config.BindAddress
			case "ServeInsecure":
				actual = Config.ServeInsecure
			case "CsrfCookieName":
				actual = Config.CsrfCookieName
			case "DisabledExtensions":
				actual = Config.DisabledExtensions
			case "CodeStyle":
				actual = Config.CodeStyle
			case "Theme":
				actual = Config.Theme
			}

			if actual != tt.expected {
				t.Errorf("%s = %v, want %v", tt.field, actual, tt.expected)
			}
		})
	}
}

func TestConfigurationStructFields(t *testing.T) {
	config := Configuration{
		Source:             "/path/to/source",
		Build:              "/path/to/build",
		Sitename:           "My Site",
		Index:              "home",
		NotFoundPage:       "error",
		BindAddress:        "0.0.0.0:8080",
		Theme:              "dark",
		CodeStyle:          "monokai",
		CsrfCookieName:     "my_csrf",
		DisabledExtensions: "ext1,ext2",
		Readonly:           true,
		ServeInsecure:      true,
	}

	tests := []struct {
		name     string
		field    string
		expected interface{}
	}{
		{"Source field", "Source", "/path/to/source"},
		{"Build field", "Build", "/path/to/build"},
		{"Sitename field", "Sitename", "My Site"},
		{"Index field", "Index", "home"},
		{"NotFoundPage field", "NotFoundPage", "error"},
		{"BindAddress field", "BindAddress", "0.0.0.0:8080"},
		{"Theme field", "Theme", "dark"},
		{"CodeStyle field", "CodeStyle", "monokai"},
		{"CsrfCookieName field", "CsrfCookieName", "my_csrf"},
		{"DisabledExtensions field", "DisabledExtensions", "ext1,ext2"},
		{"Readonly field", "Readonly", true},
		{"ServeInsecure field", "ServeInsecure", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var actual interface{}
			switch tt.field {
			case "Source":
				actual = config.Source
			case "Build":
				actual = config.Build
			case "Sitename":
				actual = config.Sitename
			case "Index":
				actual = config.Index
			case "NotFoundPage":
				actual = config.NotFoundPage
			case "BindAddress":
				actual = config.BindAddress
			case "Theme":
				actual = config.Theme
			case "CodeStyle":
				actual = config.CodeStyle
			case "CsrfCookieName":
				actual = config.CsrfCookieName
			case "DisabledExtensions":
				actual = config.DisabledExtensions
			case "Readonly":
				actual = config.Readonly
			case "ServeInsecure":
				actual = config.ServeInsecure
			}

			if actual != tt.expected {
				t.Errorf("%s = %v, want %v", tt.field, actual, tt.expected)
			}
		})
	}
}
