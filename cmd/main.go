package main

import (
	"os"

	"github.com/chai-mi/srs/internal/compile"
	domainlist "github.com/chai-mi/srs/internal/domain-list"
	"github.com/chai-mi/srs/internal/source"
	mapset "github.com/deckarep/golang-set/v2"

	"github.com/sagernet/sing-box/common/srs"
	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func main() {
	v2ray, _ := source.NewGeoSite("https://raw.githubusercontent.com/v2fly/domain-list-community/release/dlc.dat").Load()

	extra := domainlist.NewDomainList()
	extra.Suffix = map[string]mapset.Set[string]{
		"connectivitycheck.gstatic.com": mapset.NewSet("google"),
		"mcdn.bilivideo.cn":             mapset.NewSet("include-block"),
		"argotunnel.com":                mapset.NewSet("exclude-proxy"),
		"ieee.org":                      mapset.NewSet("exclude-proxy"),
		"typst.app":                     mapset.NewSet("include-proxy"),
		"chunkbase.com":                 mapset.NewSet("include-proxy"),
		"pola.rs":                       mapset.NewSet("include-proxy"),
		"neoforged.net":                 mapset.NewSet("include-proxy"),
		"tinygo.org":                    mapset.NewSet("include-proxy"),
		"bangumi.moe":                   mapset.NewSet("include-proxy"),
		"acg.rip":                       mapset.NewSet("include-proxy"),
		"share.acgnx.se":                mapset.NewSet("include-proxy"),
		"marimo.io":                     mapset.NewSet("include-proxy"),
		"qbittorrent.org":               mapset.NewSet("include-proxy"),
		"steamcontent.com":              mapset.NewSet("category-game-platforms-download"),
	}

	v2ray.Union(extra)
	v2ray.Save("./list/all.json")

	domainType := []string{"suffix", "full"}
	block := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int{
			"include-block":    1e8,
			"exclude-block":    -1e8,
			"category-ads-all": 1,
			"@ads":             1,
		},
	})
	direct := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int{
			"include-proxy":                    -1e8,
			"exclude-proxy":                    1e8,
			"category-game-platforms-download": 10,
			"geolocation-cn":                   1,
			"geolocation-!cn":                  -1,
			"category-games-cn":                1,
			"category-games-!cn":               -1,
			"@cn":                              10,
			"@!cn":                             -10,
			"category-ai-!cn":                  -1,
			"google":                           -100,
			"connectivity-check":               10,
		},
	})

	proxy := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int{
			"include-proxy":                    1e8,
			"exclude-proxy":                    -1e8,
			"category-game-platforms-download": -10,
			"geolocation-cn":                   -1,
			"geolocation-!cn":                  1,
			"category-games-cn":                -1,
			"category-games-!cn":               1,
			"@cn":                              -10,
			"@!cn":                             10,
			"category-ai-!cn":                  1,
			"google":                           100,
			"connectivity-check":               -10,
		},
	})
	ai := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int{
			"category-ai-!cn": 1,
		},
	})
	telegramGeosite := v2ray.ApplyRule(&domainlist.Rule{
		DomainType: domainType,
		TagWeight: map[string]int{
			"telegram": 1,
		},
	})

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
