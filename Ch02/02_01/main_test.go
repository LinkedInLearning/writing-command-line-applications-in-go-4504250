package main

import (
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseExe(t *testing.T) {
	exe := path.Join(t.TempDir(), "logs")
	output, err := exec.Command("go", "build", "-o", exe).CombinedOutput()
	require.NoErrorf(t, err, "build:\n%s", string(output))

	output, err = exec.Command(exe, "parse", "testdata/httpd.log").CombinedOutput()
	require.NoError(t, err, "run:\n%s", string(output))
}
