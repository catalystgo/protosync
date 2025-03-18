package config

import (
	"errors"
	"fmt"

	"github.com/catalystgo/protosync/internal/domain"
)

var (
	// Config errors
	dependencyInvalid = func(idx int) error {
		return fmt.Errorf("dependency at index: %d is invalid\n\t", idx)
	}

	ErrConfigEmpty        = fmt.Errorf("config is empty")
	ErrDirectoryEmpty     = fmt.Errorf("direcotory cannot be empty")
	ErrPathNotUnderOutDir = func(idx int, p string, outDir string) error {
		return fmt.Errorf(dependencyInvalid(idx).Error()+"Path (%s) is not under output directory: (%s). DON'T use \"..\" in path variable.", p, outDir)
	}

	// Source errors

	ErrSourceInvalid = func(idx int, source string, err error) error {
		return fmt.Errorf(dependencyInvalid(idx).Error()+"source: %s is invalid => %v", source, err)
	}
	ErrSourceAndSourcesSet = func(idx int) error {
		return errors.New(dependencyInvalid(idx).Error() + "source and sources cannot be set at the same time in dependency")
	}
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
