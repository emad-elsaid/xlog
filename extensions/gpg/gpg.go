package gpg

import (
	"flag"

	"github.com/emad-elsaid/xlog"
)

const EXT = ".md.pgp"

var gpgId string

func init() {
	flag.StringVar(&gpgId, "gpg", "", "PGP key ID to decrypt and edit .md.pgp files using gpg. if empty encryption will be off")
	app := xlog.GetApp()
	app.RegisterExtension(PGP{})
}

type PGP struct{}

func (PGP) Name() string { return "pgp" }
func (PGP) Init(app *xlog.App) {
	app.RegisterPageSource(new(encryptedPages))

	if !app.GetConfig().Readonly {
		app.RegisterCommand(commands)
		app.RequireHTMX()
		app.Post(`/+/gpg/encrypt/{page...}`, encryptHandler)
		app.Post(`/+/gpg/decrypt/{page...}`, decryptHandler)
	}
}
