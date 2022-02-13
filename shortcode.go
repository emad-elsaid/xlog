package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

func processShortCodes(content string) string {
	codes := shortCodes()
	for _, code := range codes {
		re := regexp.MustCompile(fmt.Sprintf(`(?s)\{%s\}.*?\{/%s\}`, code, code))

		start := fmt.Sprintf("{%s}", code)
		end := fmt.Sprintf("{/%s}", code)

		content = re.ReplaceAllStringFunc(content, func(match string) string {
			input := match[len(start) : len(match)-len(end)]
			return shortCode(code, input)
		})
	}

	return content
}

func shortCodes() (codes []string) {
	files, err := ioutil.ReadDir("shortcodes")
	if err != nil {
		return
	}

	for _, f := range files {
		if !f.IsDir() {
			codes = append(codes, f.Name())
		}
	}

	return
}

func shortCode(name, content string) string {
	cmd := exec.Command("shortcodes/"+name, "")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err.Error()
	}

	_, err = fmt.Fprint(stdin, content)
	if err != nil {
		return err.Error()
	}

	err = stdin.Close()
	if err != nil {
		return err.Error()
	}

	output, err := cmd.Output()
	if err != nil {
		return err.Error()
	}

	lines := strings.Split(string(output), "\n")
	mime := lines[0]
	out := strings.Join(lines[1:], "\n")

	switch mime {
	case "text/markdown":
		return renderMarkdown(out)
	case "text/html":
		return out
	default:
		return string(output)
	}
}
