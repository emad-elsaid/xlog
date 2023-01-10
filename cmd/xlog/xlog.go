package main

import (
	// Core
	"github.com/emad-elsaid/xlog"

	// All official extensions
	_ "github.com/emad-elsaid/xlog/extensions/all"
)

func main() {
	xlog.Start()
}
