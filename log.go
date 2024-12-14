package xlog

import (
	"log/slog"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"gitlab.com/greyxor/slogor"
)

func SetupLogger() {
	level := slogor.SetLevel(slog.LevelDebug)
	timeFmt := slogor.SetTimeFormat(time.TimeOnly)
	handler := slogor.NewHandler(os.Stderr, level, timeFmt)
	logger := slog.New(handler)
	slog.SetDefault(logger)

}

func callerName(f any) string {
	const xlogPrefix = "emad-elsaid/xlog/"
	const ghPrefix = "github.com/"
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	name = strings.TrimPrefix(name, ghPrefix)
	name = strings.TrimPrefix(name, xlogPrefix)
	return name
}
