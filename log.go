package xlog

import (
	"log"
	"time"
)

const (
	DEBUG = "\033[97;42m"
	INFO  = "\033[97;42m"
)

func Log(level, label, text string, args ...interface{}) func() {
	start := time.Now()
	return func() {
		if len(args) > 0 {
			log.Printf("%s %s \033[0m (%s) %s %v", level, label, time.Since(start), text, args)
		} else {
			log.Printf("%s %s \033[0m (%s) %s", level, label, time.Since(start), text)
		}
	}
}
