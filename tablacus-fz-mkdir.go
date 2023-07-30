package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/ktr0731/go-fuzzyfinder"
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
	names := []string{"plain", "proofed", "send_to_author", "proofed_by_author", "send_to_printshop"}
	idx, err := fuzzyfinder.Find(names, func(i int) string {
		return names[i]
	})
	if err != nil {
		if err == fuzzyfinder.ErrAbort {
			return 0
		}
		return 1
	}
	name := names[idx]
	var pre string
	if ds := getIndexedDirs(c); 0 < len(ds) {
		sorted := sort.StringSlice(ds)
		last := sorted[len(sorted)-1]
		dn := DirName{Name: last}
		idx := dn.incrementIndex()
		pre = fmt.Sprintf("%v_", idx)
	} else {
		pre = "0_"
	}
	newpath := filepath.Join(c, pre+name)
	if err := newDir(newpath); err != nil {
		return 1
	}
	return 0
}

type DirName struct {
	Name string
}

func (dn DirName) getIndex() string {
	idx := ""
	for i := 0; i < len(dn.Name); i++ {
		c := string(dn.Name[i])
		_, err := strconv.Atoi(c)
		if err != nil {
			break
		}
		idx += c
	}
	return idx
}

func (dn DirName) getIndexValue() int {
	idx := dn.getIndex()
	if len(idx) < 1 {
		return -1
	}
	i, err := strconv.Atoi(idx)
	if err != nil {
		return -1
	}
	return i
}

func (dn DirName) incrementIndex() string {
	i := dn.getIndexValue()
	if i < 0 {
		return ""
	}
	i += 1
	w := len(dn.getIndex())
	return fmt.Sprintf("%0"+strconv.Itoa(w)+"d", i)
}

func getIndexedDirs(root string) []string {
	var ds []string
	items, err := os.ReadDir(root)
	if err != nil {
		return ds
	}
	for _, item := range items {
		if item.IsDir() {
			n := item.Name()
			dn := DirName{Name: n}
			if -1 < dn.getIndexValue() {
				ds = append(ds, n)
			}
		}
	}
	return ds
}

func newDir(path string) error {
	s, err := os.Stat(path)
	if err == nil && s.IsDir() {
		return err
	}
	return os.Mkdir(path, os.ModePerm)
}
