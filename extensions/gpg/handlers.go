package gpg

import (
	"errors"

	"github.com/emad-elsaid/xlog"
)

var (
	deleteFailedErr     = errors.New("Couldn't delete original page")
	encryptionFailedErr = errors.New("Couldn't encrypt page")
)

func encryptHandler(w xlog.Response, r xlog.Request) xlog.Output {
	p := xlog.NewPage(r.PathValue("page"))
	if !p.Exists() {
		return xlog.NotFound("page not found")
	}

	encryptedPage := page{name: p.Name()}
	if !encryptedPage.Write(p.Content()) {
		return xlog.InternalServerError(encryptionFailedErr)
	}

	if !p.Delete() {
		return xlog.InternalServerError(deleteFailedErr)
	}

	return xlog.NoContent()
}

func decryptHandler(w xlog.Response, r xlog.Request) xlog.Output {
	p := xlog.NewPage(r.PathValue("page"))
	if !p.Exists() {
		return xlog.NotFound("page not found")
	}

	content := p.Content()
	if !p.Delete() {
		return xlog.InternalServerError(deleteFailedErr)
	}

	decryptedPage := xlog.NewPage(p.Name())
	if !decryptedPage.Write(content) {
		return xlog.InternalServerError(encryptionFailedErr)
	}

	return xlog.NoContent()
}
