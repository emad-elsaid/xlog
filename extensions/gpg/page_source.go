package gpg

import (
	"context"
	"errors"
	"io/fs"
	"path"
	"path/filepath"

	"github.com/emad-elsaid/xlog"
)

type encryptedPages struct{}

func (p *encryptedPages) Page(name string) xlog.Page {
	if len(gpgId) == 0 {
		return nil
	}

	pg := page{
		name: name,
	}
	if pg.Exists() {
		return &pg
	}

	return nil
}

func (p *encryptedPages) Each(ctx context.Context, f func(xlog.Page)) {
	if len(gpgId) == 0 {
		return
	}

	filepath.WalkDir(".", func(name string, d fs.DirEntry, err error) error {
		select {

		case <-ctx.Done():
			return errors.New("context stopped")

		default:
			lastExt := path.Ext(name)
			basename := name[:len(name)-len(lastExt)]
			secondExt := path.Ext(basename)
			ext := secondExt + lastExt
			basename = name[:len(name)-len(ext)]

			if EXT == ext {
				f(&page{
					name: basename,
				})
				break
			}

		}

		return nil
	})
}
