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

// RegisterExtension registers a new extension
func (app *App) RegisterExtension(e Extension) {
	app.extensions = append(app.extensions, e)
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
	for i := range app.extensions {
		if slices.Contains(disabled, app.extensions[i].Name()) {
			disabledNames = append(disabledNames, app.extensions[i].Name())
			continue
		}

		app.extensions[i].Init(app)
		enabledNames = append(enabledNames, app.extensions[i].Name())
	}

	slog.Info("extensions", "enabled", enabledNames, "disabled", disabled)
}
