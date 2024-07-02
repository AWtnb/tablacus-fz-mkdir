package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/zymd/dir"
)

func main() {
	var (
		cur    string
		stdout bool
	)
	flag.StringVar(&cur, "cur", "", "current dir path")
	flag.BoolVar(&stdout, "stdout", false, "switch to write new directory to stdout, instead of making directory")
	flag.Parse()
	os.Exit(run(cur, stdout))
}

func findYaml(curPath string) (string, error) {
	elems := strings.Split(curPath, string(os.PathSeparator))
	for i := 0; i <= len(elems); i++ {
		ln := len(elems) - i
		d := strings.Join(elems[0:ln], string(os.PathSeparator))
		p := filepath.Join(d, "dirnames.yaml")
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("yaml not found")
}

func newDirName(cur string) (string, bool, error) {
	yp, err := findYaml(cur)
	if err != nil {
		return "", false, err
	}
	var dns dir.DirNames
	if err := dns.Load(yp); err != nil {
		return "", false, err
	}
	return dns.Select()
}

func run(c string, stdout bool) int {

	dn, inc, err := newDirName(c)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	wd := dir.WorkDir{Path: c}
	if err := wd.Scan(); err != nil {
		fmt.Println(err)
		return 1
	}

	n := wd.WithIndex(dn)

	if stdout {
		if inc {
			fmt.Print(n)
		} else {
			fmt.Print(dn)
		}
		return 0
	}

	if !inc {
		if err := wd.Mkdir(dn); err != nil {
			return 1
		}
		return 0
	}
	if err := wd.Scan(); err != nil {
		fmt.Println(err)
		return 1
	}
	if err := wd.Mkdir(n); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
