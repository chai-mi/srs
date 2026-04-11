package source

import (
	"bufio"
	"bytes"
	"strings"

	domainlist "github.com/chai-mi/srs/internal/domain-list"
)

var _ SourceData = (*Dnsmasq)(nil)

type Dnsmasq struct {
	src     string
	addTags domainlist.Tags
}

func NewDnsmasq(src string, addTag string) *Dnsmasq {
	return &Dnsmasq{
		src:     src,
		addTags: domainlist.NewTags(addTag),
	}
}

func (dm *Dnsmasq) Load() (*domainlist.DomainList, error) {
	data, err := getInput(dm.src)
	if err != nil {
		return nil, err
	}
	return parseDnsmasq(data, dm.addTags)
}

func parseDnsmasq(input []byte, tags domainlist.Tags) (*domainlist.DomainList, error) {
	list := domainlist.NewDomainList()
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = removeServer(line)
		if len(line) == 0 {
			continue
		}
		list.Add(line, domainlist.DomainFull, tags)
	}
	return list, nil
}

func removeServer(line string) string {
	s := strings.Split(line, "/")
	if len(s) != 3 {
		return ""
	}
	return s[1]
}
