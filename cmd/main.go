package main

import (
	"os"

	"github.com/chai-mi/srs/compile"
	domainlist "github.com/chai-mi/srs/domain-list"
	"github.com/chai-mi/srs/source"

	"github.com/sagernet/sing-box/common/srs"
	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func main() {
	v2ray, _ := source.NewGeoSite("https://raw.githubusercontent.com/v2fly/domain-list-community/release/dlc.dat").Load()
	chinaList, _ := source.NewDnsmasq("https://raw.githubusercontent.com/felixonmars/dnsmasq-china-list/master/accelerated-domains.china.conf", "dnsmasq-china-list").Load()
	trackerslist, _ := source.NewUrl("https://raw.githubusercontent.com/ngosang/trackerslist/master/trackers_all.txt", "public-tracker").Load()
	trackersListCollection, _ := source.NewUrl("https://raw.githubusercontent.com/XIU2/TrackersListCollection/master/all.txt", "public-tracker").Load()
	windowsSpyBlocker, _ := source.NewHosts("https://raw.githubusercontent.com/crazy-max/WindowsSpyBlocker/master/data/hosts/spy.txt", "block-windows-spy").Load()
	extra, _ := domainlist.LoadDomainList("extra.json")

	v2ray.Save("./list/v2ray.json")

	v2ray.Union(chinaList)
	v2ray.Union(trackerslist)
	v2ray.Union(trackersListCollection)
	v2ray.Union(windowsSpyBlocker)
	v2ray.Union(extra)

	domainType := []string{"suffix", "full"}
	tracker := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int64{
			"include-tracker": 1e8,
			"exclude-tracker": -1e8,
			"public-tracker":  1,
		},
	})
	block := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int64{
			"include-block":     1e8,
			"exclude-block":     -1e8,
			"block-windows-spy": 1,
			"category-ads-all":  1,
			"@ads":              1,
		},
	})
	direct := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int64{
			"include-proxy":                    -1e8,
			"exclude-proxy":                    1e8,
			"category-game-platforms-download": 10,
			"geolocation-cn":                   1,
			"geolocation-!cn":                  -1,
			"category-games-cn":                1,
			"category-games-!cn":               -1,
			"@cn":                              10,
			"@!cn":                             -10,
			"dnsmasq-china-list":               1,
			"cn":                               1,
			"category-ai-!cn":                  -1,
			"connectivity-check":               1,
			"google":                           -100,
		},
	})

	proxy := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int64{
			"include-proxy":                    1e8,
			"exclude-proxy":                    -1e8,
			"category-game-platforms-download": -10,
			"geolocation-cn":                   -1,
			"geolocation-!cn":                  1,
			"category-games-cn":                -1,
			"category-games-!cn":               1,
			"@cn":                              -10,
			"@!cn":                             10,
			"dnsmasq-china-list":               -1,
			"cn":                               -1,
			"category-ai-!cn":                  1,
			"connectivity-check":               -1,
			"google":                           100,
		},
	})
	ai := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int64{
			"category-ai-!cn": 1,
		},
	})
	telegramGeosite := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int64{
			"telegram": 1,
		},
	})

	compile.Save2ruleset(tracker, "./public/tracker.srs")
	compile.Save2ruleset(block, "./public/block.srs")
	compile.Save2ruleset(direct, "./public/direct.srs")
	compile.Save2ruleset(proxy, "./public/proxy.srs")
	compile.Save2ruleset(ai, "./public/ai.srs")
	compile.Save2ruleset(telegramGeosite, "./public/telegram-geosite.srs")

	block.Save("./list/block.json")
	direct.Save("./list/direct.json")
	proxy.Save("./list/proxy.json")

	telegramgeoip("./public/telegram-geoip.srs")
}

func telegramgeoip(path string) error {
	var headlessRule option.DefaultHeadlessRule
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
	var plainRuleSet option.PlainRuleSet
	plainRuleSet.Rules = []option.HeadlessRule{
		{
			Type:           "default",
			DefaultOptions: headlessRule,
		},
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	srs.Write(file, plainRuleSet, constant.RuleSetVersion3)
	return nil
}
