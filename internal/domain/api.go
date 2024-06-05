package domain

import (
	"slices"
)

const (
	DefaultDomainGithub    = "github.com"
	DefaultDomainGitlab    = "gitlab.com"
	DefaultDomainBitbucket = "bitbucket.org"
)

var (
	AllowedDomainsAPI = [...]string{DefaultDomainGithub, DefaultDomainGitlab, DefaultDomainBitbucket}
)

func IsDomainAPIValid(domain string) bool {
	return slices.Contains(AllowedDomainsAPI[:], domain)
}

func GetAPIDomain(api string) string {
	switch api {
	case "github":
		return DefaultDomainGithub
	case "gitlab":
		return DefaultDomainGitlab
	case "bitbucket":
		return DefaultDomainBitbucket
	}
	return ""
}
