package xlog

type Extension interface {
	Name() string
	Init()
}

func RegisterExtension(e Extension) {
	app := GetApp()
	app.RegisterExtension(e)
}

func initExtensions() {
	app := GetApp()
	app.initExtensions()
}
