package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func hasNumberPrefix(name string) bool {
	for i := 0; i < 10; i++ {
		c := fmt.Sprint(i)
		if strings.HasPrefix(name, c) {
			return true
		}
	}
	return false
}

func getNumberPrefix(name string) string {
	pre := ""
	for i := 0; i < len(name); i++ {
		c := string(name[i])
		if _, err := strconv.Atoi(c); err != nil {
			break
		}
		pre += c
	}
	return pre
}

type WorkDir struct {
	Path    string
	subdirs []string
}

func (wd *WorkDir) Scan() error {
	items, err := os.ReadDir(wd.Path)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.IsDir() {
			wd.subdirs = append(wd.subdirs, item.Name())
		}
	}
	return nil
}

func (wd WorkDir) getIndexedMember() (dirnames []string) {
	dirs := wd.subdirs
	for _, dir := range dirs {
		if hasNumberPrefix(dir) {
			dirnames = append(dirnames, dir)
		}
	}
	return
}

func (wd WorkDir) getNameWithLargestIndex() string {
	ds := wd.getIndexedMember()
	if len(ds) < 1 {
		return ""
	}
	sorted := sort.StringSlice(ds)
	return sorted[len(sorted)-1]
}

func (wd WorkDir) WithIndex(name string) (newname string) {
	newname = fmt.Sprintf("0_%s", name)
	l := wd.getNameWithLargestIndex()
	if len(l) < 1 {
		return
	}
	p := Prefix{text: getNumberPrefix(l)}
	pre, err := p.Increment()
	if err != nil {
		return
	}
	newname = fmt.Sprintf("%s_%s", pre, name)
	return
}

func (wd WorkDir) NewDir(name string) error {
	path := filepath.Join(wd.Path, name)
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		return err
	}
	return os.Mkdir(path, os.ModePerm)
}
