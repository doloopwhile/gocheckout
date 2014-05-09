package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type vcsCmd struct {
	checkout []string
	update   []string
}

var (
	hg = &vcsCmd{
		[]string{"hg", "update"},
		[]string{"hg", "pull"},
	}
	git = &vcsCmd{
		[]string{"git", "checkout"},
		[]string{"git", "fetch"},
	}
	bzr = &vcsCmd{
		[]string{"bzr", "revert", "-r"},
		[]string{"bzr", "pull"},
	}
)

var verbose bool = false

func (vcs *vcsCmd) Checkout(dir string, revision string) error {
	args := append(vcs.checkout, revision)
	return execIn(dir, args...)
}

func (vcs *vcsCmd) Update(dir string) error {
	return execIn(dir, vcs.update...)
}

func execIn(dir string, args ...string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)
	cmd := exec.Command(args[0], args[1:]...)
	if verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func allPackageDirs(gopathDirs []string, packageName string) []string {
	var paths []string

	for _, gopathDir := range gopathDirs {
		path := filepath.Join(gopathDir, "src")
		for _, elem := range strings.Split(packageName, "/") {
			path = filepath.Join(path, elem)
			paths = append(paths, path)
		}
	}

	return paths
}

func checkout(packageName string, revision string) error {
	gopath := os.Getenv("GOPATH")
	var gopathDirs []string
	for _, dir := range strings.Split(gopath, string(filepath.ListSeparator)) {
		if len(dir) > 0 {
			gopathDirs = append(gopathDirs, dir)
		}
	}
	if len(gopathDirs) == 0 {
		return errors.New("GOPATH is empty")
	}

	for _, p := range allPackageDirs(gopathDirs, packageName) {
		var vcs *vcsCmd
		if isDir(filepath.Join(p, ".git")) {
			vcs = git
		} else if isDir(filepath.Join(p, ".hg")) {
			vcs = hg
		} else if isDir(filepath.Join(p, ".bzr")) {
			vcs = bzr
		}

		if vcs != nil {
			if verbose {
				println("Update...")
			}
			if err := vcs.Update(p); err != nil {
				return err
			}
			if verbose {
				println("Checkout...")
			}
			return vcs.Checkout(p, revision)
		}
	}

	return errors.New("Package repository not found")
}

func isDir(p string) bool {
	if fi, err := os.Stat(filepath.Join(p)); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func usage() {
	fmt.Printf("usage: %s [options] <package-name> <revision>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Usage = usage
	v := flag.Bool("v", true, "Do not omit VCS outputs")

	flag.Parse()
	if flag.NArg() != 2 {
		usage()
	}
	if *v {
		verbose = true
	}

	err := checkout(flag.Arg(0), flag.Arg(1))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("%s: ", os.Args[0]), err)
		os.Exit(1)
	}
}
