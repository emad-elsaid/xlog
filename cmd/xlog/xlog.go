package main

import (
	// Core
	"context"

	"github.com/emad-elsaid/xlog"

	// All official extensions
	_ "github.com/emad-elsaid/xlog/extensions/all"
)

func main() {
	xlog.Start(context.Background())
}
