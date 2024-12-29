package gpg

import (
	"flag"

	"github.com/emad-elsaid/xlog"
)

const EXT = ".md.pgp"

var gpgId string

func init() {
	flag.StringVar(&gpgId, "gpg", "", "PGP key ID to decrypt and edit .md.pgp files using gpg. if empty encryption will be off")
	xlog.RegisterExtension(PGP{})
}

type PGP struct{}

func (PGP) Name() string { return "pgp" }
func (PGP) Init() {
	xlog.RegisterPageSource(new(encryptedPages))

	if !xlog.Config.Readonly {
		xlog.RegisterCommand(commands)
		xlog.RequireHTMX()
		xlog.Post(`/+/gpg/encrypt/{page...}`, encryptHandler)
		xlog.Post(`/+/gpg/decrypt/{page...}`, decryptHandler)
	}
}
