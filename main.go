package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

func main() {
	var (
		cur  string
		menu string
	)
	flag.StringVar(&cur, "cur", "", "current dir path")
	flag.StringVar(&menu, "menu", "", "menu file path")
	flag.Parse()
	os.Exit(run(cur, menu))
}

func run(c string, menu string) int {
	names := getMenu(menu)
	if len(names) < 1 {
		names = []string{"plain", "proofed", "send_to_author", "proofed_by_author", "send_to_printshop", "layout"}
	}
	idx, err := fuzzyfinder.Find(names, func(i int) string {
		return names[i]
	}, fuzzyfinder.WithCursorPosition(fuzzyfinder.CursorPositionTop))
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

func readFile(path string) string {
	if !isValidPath(path) {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

func isValidPath(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func getMenu(src string) []string {
	names := []string{}
	var t string
	if isValidPath(src) {
		t = readFile(src)
	} else {
		e, _ := os.Executable()
		p := filepath.Join(filepath.Dir(e), "menu.txt")
		t = readFile(p)
	}
	if len(t) < 1 {
		return names
	}
	for _, l := range strings.Split(t, "\n") {
		l = strings.Trim(l, " ")
		l = strings.Trim(l, "\r")
		if 0 < len(l) {
			names = append(names, l)
		}
	}
	return names
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
