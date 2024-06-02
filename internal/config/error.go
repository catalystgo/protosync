package config

import "fmt"

var (
	// Config errors

	ErrConfigEmpty = fmt.Errorf("config is empty")

	// Ref errors

	ErrorInvalidRef = func(ref string, source string) error {
		return fmt.Errorf("invalid ref: %s. for source: %s", ref, source)
	}

	// Source errors

	ErrSourceEmpty       = fmt.Errorf("source cannot be empty")
	ErrSourceDomainEmpty = func(source string) error {
		return fmt.Errorf("source: %s does not have a domain", source)
	}
	ErrSourceUnregisteredDomain = func(source string, domain string) error {
		return fmt.Errorf("unregistered domain: %s for source: %s", domain, source)
	}
	ErrInvalidSource = func(source string) error {
		return fmt.Errorf("invalid source: %s. source must be like domain/user/repo/path@ref", source)
	}

	// Domain errors

	ErrDomainEmpty         = fmt.Errorf("domain cannot be empty")
	ErrDomainAlreadyExists = func(domain string) error {
		return fmt.Errorf("domain: %s already has been added", domain)
	}
	ErrInvalidDomainApi = func(domain string) error {
		return fmt.Errorf("invalid api for domain: %s. allowed domains are: %+v", domain, domainAllowedApi)
	}
)
