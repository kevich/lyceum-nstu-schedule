package tools

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		texts          []string
		expectedOutput string
		shouldOutput   bool
	}{
		{
			name:         "nil error with no text",
			err:          nil,
			texts:        []string{},
			shouldOutput: false,
		},
		{
			name:         "nil error with custom text",
			err:          nil,
			texts:        []string{"Custom error message"},
			shouldOutput: false,
		},
		{
			name:           "error with default text",
			err:            errors.New("something went wrong"),
			texts:          []string{},
			expectedOutput: "Error happened something went wrong",
			shouldOutput:   true,
		},
		{
			name:           "error with custom text",
			err:            errors.New("connection failed"),
			texts:          []string{"Failed to connect"},
			expectedOutput: "Failed to connect connection failed",
			shouldOutput:   true,
		},
		{
			name:           "error with multiple texts (only first used)",
			err:            errors.New("timeout"),
			texts:          []string{"First message", "Second message", "Third message"},
			expectedOutput: "First message timeout",
			shouldOutput:   true,
		},
		{
			name:           "error with empty string text",
			err:            errors.New("test error"),
			texts:          []string{""},
			expectedOutput: "test error",
			shouldOutput:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original stdout
			oldStdout := os.Stdout
			defer func() {
				os.Stdout = oldStdout
			}()

			// Create a pipe to capture stdout
			r, w, err := os.Pipe()
			assert.NoError(t, err, "Failed to create pipe")

			os.Stdout = w

			// Run the function
			CheckError(tt.err, tt.texts...)

			// Close write end and read output
			w.Close()
			var buf bytes.Buffer
			_, err = buf.ReadFrom(r)
			assert.NoError(t, err, "Failed to read from pipe")

			output := strings.TrimSpace(buf.String())

			if tt.shouldOutput {
				assert.Contains(t, output, tt.expectedOutput, "Output should contain expected text")
			} else {
				assert.Empty(t, output, "Should not output anything when error is nil")
			}
		})
	}
}

