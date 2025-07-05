package xlog

// RegisterProperty registers a function that returns a set of properties for a page
func (app *App) RegisterProperty(a func(Page) []Property) {
	app.propsSources = append(app.propsSources, a)
}

// Properties returns a list of properties for a page
func (app *App) Properties(p Page) map[string]Property {

	ps := map[string]Property{}
	for _, source := range app.propsSources {
		for _, pr := range source(p) {
			ps[pr.Name()] = pr
		}
	}
	return ps
}

type lastUpdateProp struct{ page Page }

func (a lastUpdateProp) Icon() string { return "fa-solid fa-clock" }
func (a lastUpdateProp) Name() string { return "modified" }
func (a lastUpdateProp) Value() any {
	app := GetApp()
	return app.ago(a.page.ModTime())
}

func DefaultProps(p Page) []Property {
	if p.ModTime().IsZero() {
		return nil
	}

	return []Property{
		lastUpdateProp{p},
	}
}
