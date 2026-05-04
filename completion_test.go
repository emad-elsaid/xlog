package xlog

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestHandleCompletion(t *testing.T) {
	tests := []struct {
		name     string
		shell    string
		contains []string
	}{
		{
			name:  "bash completion",
			shell: "bash",
			contains: []string{
				"_xlog_completion()",
				"complete -F _xlog_completion xlog",
				"-source",
				"-build",
				"-theme",
			},
		},
		{
			name:  "zsh completion",
			shell: "zsh",
			contains: []string{
				"#compdef xlog",
				"_xlog()",
				"-source[Directory that will act as a storage]",
				"-theme[bulma theme to use]:theme:(light dark)",
			},
		},
		{
			name:  "fish completion",
			shell: "fish",
			contains: []string{
				"complete -c xlog",
				"-l source -d 'Directory that will act as a storage'",
				"-l theme -d 'bulma theme to use' -r -a 'light dark'",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Mock osExit to prevent actual exit
			oldOsExit := osExit
			osExit = func(code int) {
				panic(code) // Use panic to exit the goroutine
			}
			defer func() {
				osExit = oldOsExit
			}()

			// Call handleCompletion in a goroutine to catch panic
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Expected panic from osExit(0)
						if code, ok := r.(int); !ok || code != 0 {
							t.Errorf("Expected osExit(0), got panic: %v", r)
						}
					}
				}()
				handleCompletion(tc.shell)
			}()

			// Restore stdout and read captured output
			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			// Verify expected strings are present
			for _, expected := range tc.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Output missing expected string: %q\nGot:\n%s", expected, output)
				}
			}
		})
	}
}

func TestHandleCompletionInvalidShell(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Mock osExit
	oldOsExit := osExit
	exitCode := 0
	osExit = func(code int) {
		exitCode = code
		panic(code)
	}
	defer func() {
		osExit = oldOsExit
	}()

	// Call with invalid shell
	func() {
		defer func() {
			recover() // Catch panic from osExit
		}()
		handleCompletion("invalid")
	}()

	// Restore stderr and read captured output
	w.Close()
	os.Stderr = oldStderr
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Verify error message
	if !strings.Contains(output, "Unknown shell: invalid") {
		t.Errorf("Expected error message, got: %s", output)
	}

	if !strings.Contains(output, "Supported: bash, zsh, fish") {
		t.Errorf("Expected supported shells list, got: %s", output)
	}

	// Verify exit code 1
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}
