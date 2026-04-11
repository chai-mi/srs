package source

import (
	"bufio"
	"bytes"
	"strings"

	domainlist "github.com/chai-mi/srs/internal/domain-list"
)

var _ SourceData = (*Hosts)(nil)

type Hosts struct {
	src     string
	addTags domainlist.Tags
}

func NewHosts(src string, addTag string) *Hosts {
	return &Hosts{
		src:     src,
		addTags: domainlist.NewTags(addTag),
	}
}

func (h *Hosts) Load() (*domainlist.DomainList, error) {
	data, err := getInput(h.src)
	if err != nil {
		return nil, err
	}
	return parseDnsmasq(data, h.addTags)
}

func parseHosts(input []byte, tags domainlist.Tags) (*domainlist.DomainList, error) {
	list := domainlist.NewDomainList()
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = removeComment(line)
		if len(line) == 0 {
			continue
		}
		ds := strings.Fields(line)
		if len(ds) < 2 {
			continue
		}
		for _, d := range ds[1:] {
			if d[0] == '*' {
				list.Add(d[2:], domainlist.DomainSuffix, tags)
			} else {
				list.Add(d, domainlist.DomainFull, tags)
			}
		}
	}
	return list, nil
}
