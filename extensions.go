package xlog

type Extension interface {
	Name() string
	Init()
}

var extensions = []Extension{}

func RegisterExtension(e Extension) {
	extensions = append(extensions, e)
}

func initExtensions() {
	for i := range extensions {
		extensions[i].Init()
	}
}
