package config

import (
	"github.com/catalystgo/protosync/internal/domain"
)

// validate validates the config
func validate() error {
	// configDomains is used to store domains to check if there are duplicated domains
	configDomains := make(map[string]struct{}, len(cfg.Domains)+3)
	configDomains[domain.DefaultDomainGithub] = struct{}{}
	configDomains[domain.DefaultDomainGitlab] = struct{}{}
	configDomains[domain.DefaultDomainBitbucket] = struct{}{}

	if cfg.OutDir == "" {
		return ErrOutDirEmpty
	}

	if err := validateDomain(cfg.Domains, configDomains); err != nil {
		return err
	}

	if err := validateDependencies(cfg.Dependencies, configDomains); err != nil {
		return err
	}

	return nil
}

// validateDomain validates the domain api, it checks:
// - Domain is has an allowed API
// - Domain is not duplicated
func validateDomain(domains []Domain, m map[string]struct{}) error {
	for _, d := range domains {
		// validate domain Domain
		if d.Domain == "" {
			return ErrDomainEmpty
		}

		// Validate domain API
		if !domain.IsDomainAPIValid(d.Domain) {
			return ErrInvalidDomainAPI(d.Domain)
		}

		// Check if domain has been registered before, if so return an error
		if _, ok := m[d.Domain]; ok {
			return ErrDomainAlreadyExists(d.Domain)
		}
		m[d.Domain] = struct{}{}
	}
	return nil
}

// validateDependencies validates the dependencies, it checks:
// - Source is not empty
// - Source format is valid
// - Domain is not empty
// - Domain is registered
// - Ref is not empty
func validateDependencies(deps []Dependency, configDomains map[string]struct{}) error {
	for _, dep := range deps {
		f, err := domain.ParseFile(dep.Source)
		if err != nil {
			return err
		}

		// Check if domain has been registered before or is a valid API domain if not return an error
		if _, ok := configDomains[f.Domain]; !ok || !domain.IsDomainAPIValid(f.Domain) {
			return ErrSourceUnregisteredDomain(dep.Source, f.Domain)
		}
	}

	return nil
}
