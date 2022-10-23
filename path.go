package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// lookPath search for an executable named file in the given path that is not
// in the excludedDir.
// Example:
//   - file: prog
//   - path: /Users/john/bin:/usr/local/bin
//   - both /Users/john/bin/prog and /usr/local/bin/prog exists and are
//     executable
//   - excludedDir: /Users/john/bin/
//
// this function will return /usr/local/bin/prog, not /Users/john/bin/prog.
func lookPath(file, excludedDir, path string) (string, error) {
	excludedDir = filepath.Clean(excludedDir)

	var paths []string
	for _, dir := range filepath.SplitList(path) {
		if filepath.Clean(dir) != excludedDir {
			paths = append(paths, dir)
		}
	}

	effectivePath := strings.Join(paths, string(os.PathListSeparator))

	origPath := os.Getenv("PATH")
	if err := os.Setenv("PATH", effectivePath); err != nil {
		return "", fmt.Errorf("set effective PATH: %w", err)
	}
	defer func() {
		_ = os.Setenv("PATH", origPath)
	}()
	return exec.LookPath(file)
}
