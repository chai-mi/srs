package main

import (
	"domain_generate/src/category"
	"domain_generate/src/cmd"
	"domain_generate/src/inbounds"
	"domain_generate/src/option"
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
}
