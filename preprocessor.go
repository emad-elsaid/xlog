package xlog

// A Preprocessor is a function that takes the whole page content and returns a
// modified version of the content. extensions should define this type and
// register is so that when page is rendered it will execute all of them in
// order like a pipeline each function output is passed as an input to the next.
// at the end the last preprocessor output is then rendered to HTML
type Preprocessor func(Markdown) Markdown

// RegisterPreprocessor registers a Preprocessor function. extensions should use this function to
// register a preprocessor.
func RegisterPreprocessor(f Preprocessor) {
	app := GetApp()
	app.RegisterPreprocessor(f)
}

// This function take the page content and pass it through all registered
// preprocessors and return the last preprocessor output to the caller
func PreProcess(content Markdown) Markdown {
	app := GetApp()
	return app.PreProcess(content)
}
