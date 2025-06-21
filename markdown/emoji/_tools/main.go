package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	var args []string
	cmd := "-h"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	switch cmd {
	case "github":
		githubSubCommand(args)
	case "emb-structs":
		embStructsSubCommand(args)
	case "-h":
		fallthrough
	default:
		fmt.Fprintf(os.Stderr, `Usage: _tools <subcommand> [options]
subcommands:
  github
  emb-structs
`)
		os.Exit(1)
	}
}

func usage(u func(), err error) {
	u()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	os.Exit(1)
}

func abortIfError(err error) {
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func getURL(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	abortIfError(err)
	return bs
}
