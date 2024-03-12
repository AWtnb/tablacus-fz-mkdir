package dir

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"gopkg.in/yaml.v2"
)

type DirName struct {
	Name    string
	Aliases []string
}

type DirNames struct {
	entries []DirName
}

func (dns *DirNames) Load(path string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(buf, &dns.entries); err != nil {
		return err
	}
	return nil
}

func (dns DirNames) Select() (string, error) {
	if len(dns.entries) < 1 {
		return "", fmt.Errorf("no options to pick")
	}
	idx, err := fuzzyfinder.Find(dns.entries, func(i int) string {
		e := dns.entries[i]
		if 0 < len(e.Aliases) {
			return fmt.Sprintf("%s[%s]", e.Name, strings.Join(e.Aliases, ","))
		}
		return e.Name
	})
	if err != nil {
		return "", err
	}
	de := dns.entries[idx]
	return de.Name, nil
}
