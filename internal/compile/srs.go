package compile

import (
	"os"
	"path/filepath"

	domainlist "github.com/chai-mi/srs/internal/domain-list"

	"github.com/sagernet/sing-box/common/srs"
	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func Save2ruleset(dtl *domainlist.DomainList, path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	rs := toRuleset(dtl)
	return srs.Write(file, rs, constant.RuleSetVersion3)
}

func toRuleset(dtl *domainlist.DomainList) option.PlainRuleSet {
	var headlessRule option.DefaultHeadlessRule
	if l := len(dtl.Full); l > 0 {
		headlessRule.Domain = make([]string, 0, l)
		for domain := range dtl.Full {
			headlessRule.Domain = append(headlessRule.Domain, domain)
		}
	}
	if l := len(dtl.Suffix); l > 0 {
		headlessRule.DomainSuffix = make([]string, 0, l)
		for domain := range dtl.Suffix {
			headlessRule.DomainSuffix = append(headlessRule.DomainSuffix, domain)
		}
	}
	if l := len(dtl.Keyword); l > 0 {
		headlessRule.DomainKeyword = make([]string, 0, l)
		for domain := range dtl.Keyword {
			headlessRule.DomainKeyword = append(headlessRule.DomainKeyword, domain)
		}
	}
	if l := len(dtl.Regexp); l > 0 {
		headlessRule.DomainRegex = make([]string, 0, l)
		for domain := range dtl.Regexp {
			headlessRule.DomainRegex = append(headlessRule.DomainRegex, domain)
		}
	}
	var plainRuleSet option.PlainRuleSet
	plainRuleSet.Rules = []option.HeadlessRule{
		{
			Type:           "default",
			DefaultOptions: headlessRule,
		},
	}
	return plainRuleSet
}
