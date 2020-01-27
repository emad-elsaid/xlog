package xlog

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
)

func processShortCodes(content string) string {
	codes := shortCodes()
	for _, code := range codes {
		re := regexp.MustCompile(fmt.Sprintf(`\{%s\}.*\{/%s\}`, code, code))

		start := fmt.Sprintf("{%s}", code)
		end := fmt.Sprintf("{/%s}", code)

		content = re.ReplaceAllStringFunc(content, func(match string) string {
			input := match[len(start) : len(match)-len(end)]
			return shortCode(code, input)
		})
	}

	return content
}

func shortCodes() []string {
	codes := []string{}

	files, err := ioutil.ReadDir("shortcodes")
	if err != nil {
		return codes
	}

	for _, f := range files {
		if !f.IsDir() {
			codes = append(codes, f.Name())
		}
	}

	return codes
}

func shortCode(name, content string) string {
	cmd := exec.Command("shortcodes/"+name, content)

	output, err := cmd.Output()
	if err != nil {
		return err.Error()
	}

	return string(output)
}
