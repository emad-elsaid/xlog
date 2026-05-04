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

// RegisterExtension registers a new extension. Extensions must implement the Extension
// interface and will be initialized when the server starts unless disabled via
// Config.DisabledExtensions.
func RegisterExtension(e Extension) {
	extensions = append(extensions, e)
}

func initExtensions() {
	if Config.DisabledExtensions == "all" {
		slog.Info("extensions", "disabled", "all")
		return
	}

	disabled := strings.Split(Config.DisabledExtensions, ",")
	var disabledNames []string
	enabledNames := []string{}
	for i := range extensions {
		if slices.Contains(disabled, extensions[i].Name()) {
			disabledNames = append(disabledNames, extensions[i].Name())
			continue
		}

		extensions[i].Init()
		enabledNames = append(enabledNames, extensions[i].Name())
	}

	slog.Info("extensions", "enabled", enabledNames, "disabled", disabledNames)
}
