package dir

import (
	"fmt"
	"os"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"gopkg.in/yaml.v2"
)

type DirName struct {
	Name      string
	Increment bool
	Aliases   []string
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

func (dns DirNames) Select() (string, bool, error) {
	if len(dns.entries) < 1 {
		return "", false, fmt.Errorf("no options to pick")
	}
	idx, err := fuzzyfinder.Find(dns.entries, func(i int) string {
		e := dns.entries[i]
		p := ""
		if e.Increment {
			p = "#_"
		}
		if 0 < len(e.Aliases) {
			return p + fmt.Sprintf("%s[%s]", e.Name, strings.Join(e.Aliases, ","))
		}
		return p + e.Name
	})
	if err != nil {
		return "", false, err
	}
	de := dns.entries[idx]
	return de.Name, de.Increment, nil
}
