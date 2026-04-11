package source

import (
	"strings"

	domainlist "github.com/chai-mi/srs/internal/domain-list"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

var _ SourceData = (*GeoSite)(nil)

type GeoSite struct {
	src string
}

func NewGeoSite(src string) *GeoSite {
	return &GeoSite{
		src: src,
	}
}

func (gsl *GeoSite) Load() (*domainlist.DomainList, error) {
	data, err := getInput(gsl.src)
	if err != nil {
		return nil, err
	}

	return parseGeosite(data)
}

func parseGeosite(vGeositeData []byte) (*domainlist.DomainList, error) {
	vGeositeList := routercommon.GeoSiteList{}
	err := proto.Unmarshal(vGeositeData, &vGeositeList)
	if err != nil {
		return nil, err
	}
	dtl := domainlist.NewDomainList()
	for _, vGeositeEntry := range vGeositeList.Entry {
		code := strings.ToLower(vGeositeEntry.CountryCode)
		for _, domain := range vGeositeEntry.Domain {
			tag := mapset.NewSet(code)
			if len(domain.Attribute) > 0 {
				for _, attribute := range domain.Attribute {
					tag.Add("@" + (attribute.Key))
				}
			}
			switch domain.Type {
			case routercommon.Domain_RootDomain:
				dtl.Add(domain.Value, domainlist.DomainSuffix, tag)
			case routercommon.Domain_Full:
				dtl.Add(domain.Value, domainlist.DomainFull, tag)
			case routercommon.Domain_Plain:
				dtl.Add(domain.Value, domainlist.DomainKeyword, tag)
			case routercommon.Domain_Regex:
				dtl.Add(domain.Value, domainlist.DomainRegexp, tag)
			}
		}

	}
	return dtl, nil
}
