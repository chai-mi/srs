package source

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/url"
	"strings"

	domainlist "github.com/chai-mi/srs/internal/domain-list"
)

var _ SourceData = (*Url)(nil)

type Url struct {
	src     string
	addTags domainlist.Tags
}

func NewUrl(src string, addTag string) *Url {
	return &Url{
		src:     src,
		addTags: domainlist.NewTags(addTag),
	}
}

func (u *Url) Load() (*domainlist.DomainList, error) {
	data, err := getInput(u.src)
	if err != nil {
		return nil, err
	}
	return parseDnsmasq(data, u.addTags)
}

func parseURL(input []byte, tags domainlist.Tags) (*domainlist.DomainList, error) {
	list := domainlist.NewDomainList()
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = removeComment(line)
		if len(line) == 0 {
			continue
		}
		line, err := getHost(line)
		if err != nil {
			continue
		}
		list.Add(line, domainlist.DomainFull, tags)
	}
	return list, nil
}

func getHost(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	host := u.Hostname()

	if net.ParseIP(host) != nil {
		return "", errors.New("host is not domain")
	}
	return host, nil
}
