package main

type (
	preProcessor func(string) string
)

var (
	preProcessors = []preProcessor{}
)

func PREPROCESSOR(f preProcessor) { preProcessors = append(preProcessors, f) }

func preProcess(content string) string {
	for _, v := range preProcessors {
		content = v(content)
	}

	return content
}
