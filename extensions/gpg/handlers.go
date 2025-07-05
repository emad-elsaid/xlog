package gpg

import (
	"errors"
	"net/http"

	"github.com/emad-elsaid/xlog"
)

var (
	deleteFailedErr     = errors.New("Couldn't delete original page")
	encryptionFailedErr = errors.New("Couldn't encrypt page")
)

func encryptHandler(r xlog.Request) xlog.Output {
	app := xlog.GetApp()
	p := app.NewPage(r.PathValue("page"))
	if p == nil || !p.Exists() {
		return app.NotFound("page not found")
	}

	encryptedPage := page{name: p.Name()}
	if !encryptedPage.Write(p.Content()) {
		return app.InternalServerError(encryptionFailedErr)
	}

	if !p.Delete() {
		return app.InternalServerError(deleteFailedErr)
	}

	return func(w xlog.Response, r xlog.Request) {
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
	}
}

func decryptHandler(r xlog.Request) xlog.Output {
	app := xlog.GetApp()
	p := app.NewPage(r.PathValue("page"))
	if p == nil || !p.Exists() {
		return app.NotFound("page not found")
	}

	content := p.Content()
	if !p.Delete() {
		return app.InternalServerError(deleteFailedErr)
	}

	decryptedPage := app.NewPage(p.Name())
	if decryptedPage == nil || !decryptedPage.Write(content) {
		return app.InternalServerError(encryptionFailedErr)
	}

	return func(w xlog.Response, r xlog.Request) {
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
	}
}
