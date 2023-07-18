package cmd_test

import (
	"bytes"
	"github.com/garethjevans/scm/pkg/cmd"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestTokenCmd_Run(t *testing.T) {
	tests := []struct {
		credentials string
		host        string
		expected    string
	}{
		{
			credentials: `https://user:token@github.com/`,
			host:        "https://github.com",
			expected:    "token",
		},
		{
			credentials: `https://user:token@github.com/
https://user:other@github.com/`,
			host:     "https://github.com",
			expected: "token",
		},
		{
			credentials: `https://user:token@github.com/
https://user:other@gitlab.com/`,
			host:     "https://github.com",
			expected: "token",
		},
		{
			credentials: `https://user:token@github.com/
https://user:other@gitlab.com/`,
			host:     "https://gitlab.com",
			expected: "other",
		},
	}

	u := cmd.NewTokenCmd()

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			file, err := os.CreateTemp("", "test")
			assert.NoError(t, err)

			defer os.Remove(file.Name())

			// Example writing to the file
			text := []byte(tc.credentials)
			if _, err = file.Write(text); err != nil {
				assert.NoError(t, err)
			}

			// Close the file
			if err := file.Close(); err != nil {
				assert.NoError(t, err)
			}

			cmd.Path = file.Name()
			cmd.Host = tc.host

			b := bytes.NewBufferString("")
			u.SetOut(b)

			err = u.Execute()
			assert.NoError(t, err)

			out, err := io.ReadAll(b)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(out))
		})
	}
}
