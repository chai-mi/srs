package source

import domainlist "github.com/chai-mi/srs/internal/domain-list"

type SourceData interface {
	Load() (*domainlist.DomainList, error)
}
