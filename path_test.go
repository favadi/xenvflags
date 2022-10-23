package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createExecutable(t *testing.T, file string) {
	t.Helper()

	f, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0755)
	require.NoError(t, err)
	require.NoError(t, f.Close())
}

func TestLookPath(t *testing.T) {
	const cmd = "prog"

	path1 := t.TempDir()
	path2 := t.TempDir()
	pathEnv := fmt.Sprintf("%s:%s", path1, path2)

	_, err := lookPath(cmd, "", pathEnv)
	assert.ErrorIs(t, err, exec.ErrNotFound)

	// executable is located in path2
	createExecutable(t, filepath.Join(path2, cmd))
	executable, err := lookPath(cmd, "", pathEnv)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(path2, cmd), executable)

	// executable is located in both path1 and path2
	createExecutable(t, filepath.Join(path1, cmd))
	executable, err = lookPath(cmd, "", pathEnv)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(path1, cmd), executable)

	// executable is located in both path1 and path2, path1 is excluded
	executable, err = lookPath(cmd, path1, pathEnv)
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(path2, cmd), executable)
}
