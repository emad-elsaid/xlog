package xlog

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestHandleCompletion(t *testing.T) {
	tests := []struct {
		name          string
		shell         string
		expectExit    bool
		exitCode      int
		expectStdout  string
		expectStderr  string
		checkContains bool
	}{
		{
			name:          "bash completion",
			shell:         "bash",
			expectExit:    false,
			expectStdout:  bashCompletionScript,
			expectStderr:  "",
			checkContains: true,
		},
		{
			name:          "zsh completion",
			shell:         "zsh",
			expectExit:    false,
			expectStdout:  zshCompletionScript,
			expectStderr:  "",
			checkContains: true,
		},
		{
			name:          "fish completion",
			shell:         "fish",
			expectExit:    false,
			expectStdout:  fishCompletionScript,
			expectStderr:  "",
			checkContains: true,
		},
		{
			name:          "unknown shell",
			shell:         "powershell",
			expectExit:    true,
			exitCode:      1,
			expectStdout:  "",
			expectStderr:  "Unknown shell: powershell",
			checkContains: true,
		},
		{
			name:          "empty shell",
			shell:         "",
			expectExit:    true,
			exitCode:      1,
			expectStdout:  "",
			expectStderr:  "Unknown shell: ",
			checkContains: true,
		},
		{
			name:          "unsupported shell tcsh",
			shell:         "tcsh",
			expectExit:    true,
			exitCode:      1,
			expectStdout:  "",
			expectStderr:  "Unknown shell: tcsh. Supported: bash, zsh, fish",
			checkContains: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			rOut, wOut, _ := os.Pipe()
			os.Stdout = wOut

			// Capture stderr
			oldStderr := os.Stderr
			rErr, wErr, _ := os.Pipe()
			os.Stderr = wErr

			// Mock osExit
			oldOsExit := osExit
			exitCalled := false
			exitCodeReceived := -1
			osExit = func(code int) {
				exitCalled = true
				exitCodeReceived = code
				panic("exit called") // Stop execution
			}

			// Restore all at end
			defer func() {
				os.Stdout = oldStdout
				os.Stderr = oldStderr
				osExit = oldOsExit
			}()

			// Run the function
			if tc.expectExit {
				// Expect panic from mocked osExit
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected os.Exit to be called, but it wasn't")
					}
				}()
			}

			handleCompletion(tc.shell)

			// Close writers and read output
			wOut.Close()
			wErr.Close()

			var bufOut bytes.Buffer
			var bufErr bytes.Buffer
			_, _ = bufOut.ReadFrom(rOut)
			_, _ = bufErr.ReadFrom(rErr)

			stdout := bufOut.String()
			stderr := bufErr.String()

			// Verify exit behavior
			if tc.expectExit {
				if !exitCalled {
					t.Error("Expected os.Exit to be called")
				}
				if exitCodeReceived != tc.exitCode {
					t.Errorf("Expected exit code %d, got %d", tc.exitCode, exitCodeReceived)
				}
			} else if exitCalled {
				t.Errorf("Did not expect os.Exit to be called, but got exit code %d", exitCodeReceived)
			}

			// Verify stdout
			if tc.checkContains {
				if tc.expectStdout != "" && !strings.Contains(stdout, tc.expectStdout) {
					t.Errorf("Expected stdout to contain:\n%s\nGot:\n%s",
						tc.expectStdout[:min(len(tc.expectStdout), 100)],
						stdout[:min(len(stdout), 100)])
				}
			} else {
				if stdout != tc.expectStdout {
					t.Errorf("Expected stdout:\n%s\nGot:\n%s", tc.expectStdout, stdout)
				}
			}

			// Verify stderr
			if tc.checkContains {
				if tc.expectStderr != "" && !strings.Contains(stderr, tc.expectStderr) {
					t.Errorf("Expected stderr to contain: %q, got: %q", tc.expectStderr, stderr)
				}
			} else {
				if stderr != tc.expectStderr {
					t.Errorf("Expected stderr: %q, got: %q", tc.expectStderr, stderr)
				}
			}
		})
	}
}

func TestHandleCompletion_BashScriptContent(t *testing.T) {
	// Verify bash completion script contains essential elements
	requiredElements := []string{
		"_xlog_completion",
		"COMPREPLY",
		"-bind",
		"-build",
		"-completion",
		"-source",
		"complete -F _xlog_completion xlog",
	}

	for _, element := range requiredElements {
		if !strings.Contains(bashCompletionScript, element) {
			t.Errorf("bash completion script missing required element: %s", element)
		}
	}
}

func TestHandleCompletion_ZshScriptContent(t *testing.T) {
	// Verify zsh completion script contains essential elements
	requiredElements := []string{
		"#compdef xlog",
		"_xlog",
		"-bind",
		"-build",
		"-completion",
		"-source",
		"_arguments",
	}

	for _, element := range requiredElements {
		if !strings.Contains(zshCompletionScript, element) {
			t.Errorf("zsh completion script missing required element: %s", element)
		}
	}
}

func TestHandleCompletion_FishScriptContent(t *testing.T) {
	// Verify fish completion script contains essential elements
	requiredElements := []string{
		"complete -c xlog",
		"-l bind",
		"-l build",
		"-l completion",
		"-l source",
		"bash zsh fish",
	}

	for _, element := range requiredElements {
		if !strings.Contains(fishCompletionScript, element) {
			t.Errorf("fish completion script missing required element: %s", element)
		}
	}
}
