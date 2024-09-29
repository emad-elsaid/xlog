package xlog

import (
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/golang-cz/devslog"
)

func init() {
	slog.SetDefault(slog.New(devslog.NewHandler(os.Stdout, nil)))
}

func timing(msg string, args ...any) func() {
	start := time.Now()
	l := slog.With(args...)
	return func() {
		l.Info(msg, "time", time.Since(start))
	}
}

func FuncName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
