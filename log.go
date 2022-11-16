package xlog

import (
	"log"
	"time"
)

func timing(label, text string) func() {
	start := time.Now()
	return func() {
		log.Printf("%10s %10s %s", "["+label+"]", time.Since(start), text)
	}
}
