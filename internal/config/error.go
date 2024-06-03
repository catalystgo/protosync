package config

import (
	"fmt"

	"github.com/catalystgo/protosync/internal/domain"
)

var (
	// Config errors

	ErrConfigEmpty = fmt.Errorf("config is empty")
	ErrOutDirEmpty = fmt.Errorf("out_dir cannot be empty")

	// Source errors

	ErrSourceUnregisteredDomain = func(source string, domain string) error {
		return fmt.Errorf("unregistered domain: %s for source: %s", domain, source)
	}

	// Domain errors

	ErrDomainEmpty         = fmt.Errorf("domain cannot be empty")
	ErrDomainAlreadyExists = func(domain string) error {
		return fmt.Errorf("domain: %s already has been added", domain)
	}
	ErrInvalidDomainAPI = func(d string) error {
		return fmt.Errorf("invalid api for domain: %s. allowed API domains are: %+v", d, domain.AllowedDomainsAPI)
	}
)
