package xlog

import (
	"io/fs"
)

type defaultedFS struct {
	fs       fs.FS
	fallback fs.FS
}

func (df defaultedFS) Open(name string) (fs.File, error) {
	f, err := df.fs.Open(name)
	if err == nil {
		return f, err
	}

	return df.fallback.Open(name)
}
