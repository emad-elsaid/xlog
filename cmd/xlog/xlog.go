package main

import (
	// Core
	"context"

	"github.com/emad-elsaid/xlog"

	// All official extensions
	_ "github.com/emad-elsaid/xlog/extensions/all"
)

func main() {
	app := xlog.GetApp()
	app.Start(context.Background())
}
