package main

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateExe(t *testing.T) {
	exe := buildExe(t)
	output, err := exec.Command(exe, "validate", "testdata/httpd.log").CombinedOutput()
	require.NoError(t, err, "run:\n%s", string(output))
}

func TestValidateExeStdin(t *testing.T) {
	file, err := os.Open("testdata/httpd.log")
	require.NoError(t, err, "open")
	defer file.Close()

	exe := buildExe(t)

	cmd := exec.Command(exe, "validate")
	cmd.Stdin = file

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "run:\n%s", string(output))
}

func TestValidateExeEnv(t *testing.T) {
	exe := buildExe(t)
	cmd := exec.Command(exe, "validate", "testdata/httpd.log")
	cmd.Env = append(os.Environ(), "LOGS_VERBOSE=true")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "run:\n%s", string(output))
	require.Contains(t, string(output), "successfully validated", "output")
}

func buildExe(t *testing.T) string {
	exe := path.Join(t.TempDir(), "logs")
	output, err := exec.Command("go", "build", "-o", exe).CombinedOutput()
	require.NoErrorf(t, err, "build:\n%s", string(output))

	return exe
}
