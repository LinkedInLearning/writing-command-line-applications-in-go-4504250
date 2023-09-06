package main

import (
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
Write tests that will execute rot13.
- Test with text as argument
	- Test with various inputs (table test, see t.Run)
- Test with text from standard input
*/

// https://go-proverbs.github.io/
var rot13Cases = []struct {
	text    string
	encoded string
}{
	{"Cgo is not Go.", "Ptb vf abg Tb."},
	{"Errors are values.", "Reebef ner inyhrf."},
	{"Don't panic.", "Qba'g cnavp."},
}

func TestRot13(t *testing.T) {
	exe := buildExe(t)

	for _, tc := range rot13Cases {
		t.Run(tc.text, func(t *testing.T) {
			out, err := exec.Command(exe, tc.text).CombinedOutput()
			require.NoError(t, err, "run")
			encoded := strings.TrimSpace(string(out))
			require.Equal(t, tc.encoded, encoded)
		})
	}
}

func TestRot13Stdin(t *testing.T) {
	exe := buildExe(t)

	for _, tc := range rot13Cases {
		t.Run(tc.text, func(t *testing.T) {
			cmd := exec.Command(exe)
			cmd.Stdin = strings.NewReader(tc.text)
			out, err := cmd.CombinedOutput()
			require.NoError(t, err, "run")
			encoded := strings.TrimSpace(string(out))
			require.Equal(t, tc.encoded, encoded)
		})
	}
}

func buildExe(t *testing.T) string {
	exe := path.Join(t.TempDir(), "logs")
	output, err := exec.Command("go", "build", "-o", exe).CombinedOutput()
	require.NoErrorf(t, err, "build:\n%s", string(output))

	return exe
}
