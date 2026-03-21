package source

import domainlist "github.com/chai-mi/srs/domain-list"

type SourceData interface {
	Load() (*domainlist.DomainList, error)
}
