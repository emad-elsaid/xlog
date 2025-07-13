package xlog

import (
	"log/slog"
	"slices"
	"strings"
)

// Extension represents a plugin that can be registered with the application
type Extension interface {
	Name() string
	Init(*App)
}

// Global variable to store all registered extensions
var globalExtensions []Extension

// RegisterExtension registers a new extension
func RegisterExtension(e Extension) {
	globalExtensions = append(globalExtensions, e)
}

// initExtensions initializes all registered extensions
func (app *App) initExtensions() {
	if app.config.DisabledExtensions == "all" {
		slog.Info("extensions", "disabled", "all")
		return
	}

	disabled := strings.Split(app.config.DisabledExtensions, ",")
	disabledNames := []string{} // because the user can input wrong extension name
	enabledNames := []string{}
	for i := range globalExtensions {
		if slices.Contains(disabled, globalExtensions[i].Name()) {
			disabledNames = append(disabledNames, globalExtensions[i].Name())
			continue
		}

		globalExtensions[i].Init(app)
		enabledNames = append(enabledNames, globalExtensions[i].Name())
	}

	slog.Info("extensions", "enabled", enabledNames, "disabled", disabled)
}
