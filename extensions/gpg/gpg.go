package gpg

import (
	"flag"

	"github.com/emad-elsaid/xlog"
)

const EXT = ".md.pgp"

var gpgId string

func init() {
	flag.StringVar(&gpgId, "gpg", "", "PGP key ID to decrypt and edit .md.pgp files using gpg. if empty encryption will be off")
	xlog.RegisterPageSource(new(encryptedPages))
	xlog.RegisterQuickCommand(quickCommands)
	xlog.Post(`/+/gpg/encrypt/{page...}`, encryptHandler)
	xlog.Post(`/+/gpg/decrypt/{page...}`, decryptHandler)
}
