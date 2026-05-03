package gpg

import (
	"errors"
	"net/http"

	"github.com/emad-elsaid/xlog"
)

var (
	errDeleteFailed     = errors.New("couldn't delete original page")
	errEncryptionFailed = errors.New("couldn't encrypt page")
)

func encryptHandler(r xlog.Request) xlog.Output {
	p := xlog.NewPage(r.PathValue("page"))
	if p == nil || !p.Exists() {
		return xlog.NotFound("page not found")
	}

	encryptedPage := page{name: p.Name()}
	if !encryptedPage.Write(p.Content()) {
		return xlog.InternalServerError(errEncryptionFailed)
	}

	if !p.Delete() {
		return xlog.InternalServerError(errDeleteFailed)
	}

	return func(w xlog.Response, r xlog.Request) {
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
	}
}

func decryptHandler(r xlog.Request) xlog.Output {
	p := xlog.NewPage(r.PathValue("page"))
	if p == nil || !p.Exists() {
		return xlog.NotFound("page not found")
	}

	content := p.Content()
	if !p.Delete() {
		return xlog.InternalServerError(errDeleteFailed)
	}

	decryptedPage := xlog.NewPage(p.Name())
	if decryptedPage == nil || !decryptedPage.Write(content) {
		return xlog.InternalServerError(errEncryptionFailed)
	}

	return func(w xlog.Response, r xlog.Request) {
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
	}
}
