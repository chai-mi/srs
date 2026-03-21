package domainlist

import (
	"encoding/json"
	"os"
	"path/filepath"

	mapset "github.com/deckarep/golang-set/v2"
)

type DomainList struct {
	Full    DomainTags `json:"full,omitempty"`
	Suffix  DomainTags `json:"suffix,omitempty"`
	Keyword DomainTags `json:"keyword,omitempty"`
	Regexp  DomainTags `json:"regexp,omitempty"`
}

type DomainTags map[string]Tags

type Tags = mapset.Set[string]

func NewTags(tag string) Tags {
	return mapset.NewSet(tag)
}

const (
	DomainFull    = "full"
	DomainSuffix  = "suffix"
	DomainKeyword = "keyword"
	DomainRegexp  = "regexp"
)

func NewDomainList() *DomainList {
	return &DomainList{
		Full:    make(DomainTags),
		Suffix:  make(DomainTags),
		Keyword: make(DomainTags),
		Regexp:  make(DomainTags),
	}
}

func LoadDomainList(path string) (*DomainList, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var dl DomainList

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dl)
	if err != nil {
		return nil, err
	}
	return &dl, nil
}

func (dl *DomainList) Save(path string) error {
	jsonData, err := json.MarshalIndent(dl, "", "    ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (dl *DomainList) Union(otherDomainList *DomainList) {
	if otherDomainList == nil {
		return
	}
	for domain, tags := range otherDomainList.Full {
		if dl.Full[domain] == nil {
			dl.Full[domain] = tags
		} else {
			dl.Full[domain] = dl.Full[domain].Union(tags)
		}
	}
	for domain, tags := range otherDomainList.Suffix {
		if dl.Suffix[domain] == nil {
			dl.Suffix[domain] = tags
		} else {
			dl.Suffix[domain] = dl.Suffix[domain].Union(tags)
		}
	}
	for domain, tags := range otherDomainList.Keyword {
		if dl.Keyword[domain] == nil {
			dl.Keyword[domain] = tags
		} else {
			dl.Keyword[domain] = dl.Keyword[domain].Union(tags)
		}
	}
	for domain, tags := range otherDomainList.Regexp {
		if dl.Regexp[domain] == nil {
			dl.Regexp[domain] = tags
		} else {
			dl.Regexp[domain] = dl.Regexp[domain].Union(tags)
		}
	}
}

func (dl *DomainList) Add(domain string, domainType string, tags Tags) {
	if tags == nil {
		return
	}
	switch domainType {
	case DomainFull:
		if dl.Full[domain] == nil {
			dl.Full[domain] = tags
		} else {
			dl.Full[domain] = dl.Full[domain].Union(tags)
		}
	case DomainSuffix:
		if dl.Suffix[domain] == nil {
			dl.Suffix[domain] = tags
		} else {
			dl.Suffix[domain] = dl.Suffix[domain].Union(tags)
		}
	case DomainKeyword:
		if dl.Keyword[domain] == nil {
			dl.Keyword[domain] = tags
		} else {
			dl.Keyword[domain] = dl.Keyword[domain].Union(tags)
		}
	case DomainRegexp:
		if dl.Regexp[domain] == nil {
			dl.Regexp[domain] = tags
		} else {
			dl.Regexp[domain] = dl.Regexp[domain].Union(tags)
		}
	}
}

func (dl *DomainList) AddTag(tags Tags) {
	if tags == nil {
		return
	}
	for domain := range dl.Full {
		if dl.Full[domain] == nil {
			dl.Full[domain] = tags
		} else {
			dl.Full[domain].Union(tags)
		}
	}
	for domain := range dl.Suffix {
		if dl.Suffix[domain] == nil {
			dl.Suffix[domain] = tags
		} else {
			dl.Suffix[domain].Union(tags)
		}
	}
	for domain := range dl.Keyword {
		if dl.Keyword[domain] == nil {
			dl.Keyword[domain] = tags
		} else {
			dl.Keyword[domain].Union(tags)
		}
	}
	for domain := range dl.Regexp {
		if dl.Regexp[domain] == nil {
			dl.Regexp[domain] = tags
		} else {
			dl.Regexp[domain].Union(tags)
		}
	}
}
