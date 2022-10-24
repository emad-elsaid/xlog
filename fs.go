package xlog

import (
	"io/fs"
	"log"
)

// return file that exists in one of the FS structs.
// Prioritizing the end of the slice over earlier FSs.
type priorityFS []fs.FS

func (df priorityFS) Open(name string) (fs.File, error) {
	log.Println(name, "fs:", df)

	for i := len(df) - 1; i >= 0; i-- {
		cf := df[i]
		f, err := cf.Open(name)
		if err == nil {
			return f, err
		}
	}

	return nil, fs.ErrNotExist
}
