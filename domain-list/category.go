package domainlist

import "slices"

type Rule struct {
	Over       int64            `json:"over,omitempty"`
	DomainType []string         `json:"domain type,omitempty"`
	TagWeight  map[string]int64 `json:"tag weight"`
}

func (dl *DomainList) ApplyRule(rule *Rule) *DomainList {
	dr := NewDomainList()
	for domain, tags := range dl.Full {
		if slices.Contains(rule.DomainType, DomainFull) && countWeight(tags, rule.TagWeight) > rule.Over {
			dr.Full[domain] = tags
		}
	}
	for domain, tags := range dl.Suffix {
		if slices.Contains(rule.DomainType, DomainSuffix) && countWeight(tags, rule.TagWeight) > rule.Over {
			dr.Suffix[domain] = tags
		}
	}
	for domain, tags := range dl.Keyword {
		if slices.Contains(rule.DomainType, DomainKeyword) && countWeight(tags, rule.TagWeight) > rule.Over {
			dr.Keyword[domain] = tags
		}
	}
	for domain, tags := range dl.Regexp {
		if slices.Contains(rule.DomainType, DomainRegexp) && countWeight(tags, rule.TagWeight) > rule.Over {
			dr.Regexp[domain] = tags
		}
	}
	return dr
}

func countWeight(tags Tags, tagWeight map[string]int64) int64 {
	var count int64
	for tag := range tags.Iter() {
		count += tagWeight[tag]
	}
	return count
}
