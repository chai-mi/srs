package main

import (
	"domain_generate/src/category"
	"domain_generate/src/cmd"
	"domain_generate/src/inbounds"
	"domain_generate/src/option"
	"os"

	"github.com/sagernet/sing-box/common/srs"
	"github.com/sagernet/sing-box/constant"
	soption "github.com/sagernet/sing-box/option"
)

func main() {
	cmd.GetCmd()
	config, err := option.LoadConfig(*cmd.Config)
	if err != nil {
		return
	}

	domainList := inbounds.LoadAll(config.Inbounds)
	if *cmd.PutAllDomains != "" {
		domainList.Save(*cmd.PutAllDomains)
	}

	result := category.ParseRules(domainList, config.Rules)
	result.Save(config.Outbounds)

	telegramgeoip()
}

func telegramgeoip() error {
	var headlessRule soption.DefaultHeadlessRule
	headlessRule.IPCIDR = []string{
		"91.105.192.0/23",
		"91.108.4.0/22",
		"91.108.8.0/21",
		"91.108.16.0/21",
		"91.108.56.0/22",
		"149.154.160.0/20",
		"185.76.151.0/24",
		"2001:67c:4e8::/48",
		"2001:b28:f23c::/47",
		"2001:b28:f23f::/48",
		"2a0a:f280::/32",
	}
	var plainRuleSet soption.PlainRuleSet
	plainRuleSet.Rules = []soption.HeadlessRule{
		{
			Type:           "default",
			DefaultOptions: headlessRule,
		},
	}
	file, err := os.Create("public/telegram-geoip.srs")
	if err != nil {
		return err
	}
	defer file.Close()
	srs.Write(file, plainRuleSet, constant.RuleSetVersion3)
	return nil
}
