package xlog

// RegisterPreprocessor registers a preprocessor function
func (app *App) RegisterPreprocessor(f Preprocessor) {
	app.preprocessors = append(app.preprocessors, f)
}

// PreProcess processes content through all registered preprocessors
func (app *App) PreProcess(content Markdown) Markdown {

	for _, v := range app.preprocessors {
		content = v(content)
	}

	return content
}
