package xlog

import (
	"log"
	"reflect"
	"runtime"
	"time"
)

func timing(label, text string) func() {
	start := time.Now()
	return func() {
		log.Printf("[%10s] %-10s %s", label, time.Since(start), text)
	}
}

func FuncName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
