package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/shell"
)

const (
	envSuffix = "_EXTRA_ARGS"

	usage = `xenvflags should not be called directly, but via a symlink.

For example, shfmt is installed in /usr/local/bin/shfmt.
Setup: ln -s path/to/real/xenvflags $HOME/bin/shfmt

Make sure that $HOME/bin is placed before /usr/local/bin in PATH. Calling shfmt
will automatically apply extra arguments from SHFMT_EXTRA_ARGS environment
variable.

SHFMT_EXTRA_ARGS='-i 2' shfmt --> shfmt -i 2
`
)

var (
	version = "dev"
)

func isDebug() bool {
	return os.Getenv("XENVFLAGS_DEBUG") == "true"
}

func printVersion() {
	if os.Getenv("XENVFLAGS_VERSION") == "true" {
		fmt.Println(version)
		os.Exit(0)
	}
}

// isSymlink returns an error if the given file is not a symlink.
func isSymlink(file string) error {
	fileInfo, err := os.Lstat(file)
	if err != nil {
		return err
	}
	if !(fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink) {
		return fmt.Errorf("%s is not a symlink", file)
	}
	return nil
}

// findRealExecutable lookups for a `real` command in PATH with same name as
// the given symlink file.
// For example:
// - file: $HOME/bin/shfmt is a symlink to xenvflags
// - real shfmt is installed in /usr/local/bin/shfmt
// - PATH: $HOME/bin/shfmt:/usr/local/bin
// This function will return /usr/local/bin/shfmt.
func findRealExecutable(file string) (string, error) {
	cmd := filepath.Base(file)
	excludedDir := filepath.Dir(file)
	return lookPath(cmd, excludedDir)
}

func getExtraArgs(cmdName string) ([]string, error) {
	envArgs := os.Getenv(strings.ToUpper(cmdName) + envSuffix)
	extraArgs, err := shell.Fields(envArgs, nil)
	if err != nil {
		return nil, err
	}
	return extraArgs, nil
}

func main() {
	printVersion()

	symExecutable, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}

	if err = isSymlink(symExecutable); err != nil {
		if isDebug() {
			log.Print(err.Error())
		}
		fmt.Fprintf(os.Stderr, usage)
		os.Exit(1)
	}

	executable, err := findRealExecutable(symExecutable)
	if err != nil {
		log.Fatalf("find real executable: %s", err.Error())
	}

	origArgs := os.Args[1:]

	extraArgs, err := getExtraArgs(filepath.Base(executable))
	if err != nil {
		log.Fatalf("parse extra args from environment variable: %s", err.Error())
	}

	if isDebug() {
		log.Printf("version: %s", version)
		log.Printf("executable: %s", executable)
		log.Printf("original arguments: %s", origArgs)
		log.Printf("extra arguments from env: %s", extraArgs)
	}

	cmd := exec.Command(executable, append(extraArgs, origArgs...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			os.Exit(exitError.ExitCode())
		}

		log.Fatal(err.Error())
	}
}
