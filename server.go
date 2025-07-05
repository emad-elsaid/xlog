package xlog

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/gorilla/csrf"
)

var (
	// a function that returns the CSRF token
	CSRF = csrf.Token
)

type (
	// alias http.ResponseWriter for shorter handler declaration
	Response = http.ResponseWriter
	// alias *http.Request for shorter handler declaration
	Request = *http.Request
	// alias of http.HandlerFunc as output is expected from defined http handlers
	Output = http.HandlerFunc
	// map of string to any value used for template rendering
	Locals map[string]any // passed to templates
)

// HandlerFunc is the type of an HTTP handler function + returns output function.
// it makes it easier to return the output directly instead of writing the output to w then return.
type HandlerFunc func(Request) Output

type funcStringer struct{ any }

func (f funcStringer) String() string {
	const xlogPrefix = "emad-elsaid/xlog/"
	const ghPrefix = "github.com/"
	name := runtime.FuncForPC(reflect.ValueOf(f.any).Pointer()).Name()
	name = strings.TrimPrefix(name, ghPrefix)
	name = strings.TrimPrefix(name, xlogPrefix)
	return name
}
