package cmd_test

import (
	"bytes"
	"github.com/garethjevans/scm/pkg/cmd"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestUrlCmd_Run(t *testing.T) {
	tests := []struct {
		host     string
		owner    string
		repo     string
		kind     string
		expected string
	}{
		{
			host:     "https://github.com/",
			owner:    "garethjevans",
			repo:     "scm",
			kind:     "",
			expected: "https://github.com/garethjevans/scm",
		},
		{
			host:     "https://github.com",
			owner:    "garethjevans",
			repo:     "scm",
			kind:     "",
			expected: "https://github.com/garethjevans/scm",
		},
		{
			host:     "https://dev.azure.com",
			owner:    "garethjevans",
			repo:     "scm",
			kind:     "",
			expected: "https://dev.azure.com/garethjevans/_git/scm",
		},
		{
			host:     "https://dev.azure.com/",
			owner:    "garethjevans",
			repo:     "scm",
			kind:     "",
			expected: "https://dev.azure.com/garethjevans/_git/scm",
		},
		{
			host:     "https://custom.azdo.host/",
			owner:    "garethjevans",
			repo:     "scm",
			kind:     "azure",
			expected: "https://custom.azdo.host/garethjevans/_git/scm",
		},
	}

	u := cmd.NewUrlCmd()

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			cmd.Owner = tc.owner
			cmd.Host = tc.host
			cmd.Repo = tc.repo
			cmd.Kind = tc.kind

			b := bytes.NewBufferString("")
			u.SetOut(b)

			err := u.Execute()
			assert.NoError(t, err)

			out, err := io.ReadAll(b)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(out))
		})
	}
}
