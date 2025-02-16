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
	if Config.DisabledExtensions == "all" {
		slog.Info("extensions", "disabled", "all")
		return
	}

	disabled := strings.Split(Config.DisabledExtensions, ",")
	disabledNames := []string{} // because the user can input wrong extension name
	enabledNames := []string{}
	for i := range extensions {
		if slices.Contains(disabled, extensions[i].Name()) {
			disabledNames = append(disabledNames, extensions[i].Name())
			continue
		}

		extensions[i].Init()
		enabledNames = append(enabledNames, extensions[i].Name())
	}

	slog.Info("extensions", "enabled", enabledNames, "disabled", disabled)
}
