//go:build !appengine && !js && go1.21
// +build !appengine,!js,go1.21

package util

import (
	"unsafe"
)

// BytesToReadOnlyString returns a string converted from given bytes.
func BytesToReadOnlyString(b []byte) string {
	// #nosec G103 -- Intentional unsafe optimization for zero-copy string conversion (Go 1.21+)
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToReadOnlyBytes returns bytes converted from given string.
func StringToReadOnlyBytes(s string) []byte {
	// #nosec G103 -- Intentional unsafe optimization for zero-copy byte conversion (Go 1.21+)
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
