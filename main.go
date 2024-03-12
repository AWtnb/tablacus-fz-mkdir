package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AWtnb/tablacus-fz-mkdir/dir"
)

func main() {
	var (
		cur string
	)
	flag.StringVar(&cur, "cur", "", "current dir path")
	flag.Parse()
	os.Exit(run(cur))
}

func run(c string) int {
	yp, err := findYaml(c)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	var dns dir.DirNames
	if err := dns.Load(yp); err != nil {
		fmt.Println(err)
		return 1
	}

	dn, err := dns.Select()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	wd := dir.WorkDir{Path: c}
	if err := wd.Scan(); err != nil {
		fmt.Println(err)
		return 1
	}
	nn := wd.WithIndex(dn)
	if err := wd.NewDir(nn); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
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
