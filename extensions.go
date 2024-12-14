package xlog

import (
	"log/slog"
	"slices"
	"strings"
)

type Extension interface {
	Name() string
	Init()
}

var extensions = []Extension{}

func RegisterExtension(e Extension) {
	extensions = append(extensions, e)
}

func initExtensions() {
	disabled := strings.Split(Config.DisabledExtensions, ",")
	names := []string{}
	for i := range extensions {
		if slices.Contains(disabled, extensions[i].Name()) {
			continue
		}

		extensions[i].Init()
		names = append(names, extensions[i].Name())
	}

	slog.Info("extensions", "names", names)
}
