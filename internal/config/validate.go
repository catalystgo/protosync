package config

import (
	"slices"
	"strings"
)

// validate validates the config
func validate() error {
	// configDomains is used to store domains to check if there are duplicated domains
	configDomains := make(map[string]struct{}, len(cfg.Domains))

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
	for _, domain := range domains {
		// validate domain Domain
		if domain.Domain == "" {
			return ErrDomainEmpty
		}

		// Validate domain API
		if !slices.Contains(domainAllowedApi, domain.Api) {
			return ErrInvalidDomainApi(domain.Domain)
		}
		if _, ok := m[domain.Domain]; ok {
			return ErrDomainAlreadyExists(domain.Domain)
		}
		m[domain.Domain] = struct{}{}
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
		if dep.Source == "" {
			return ErrSourceEmpty
		}

		// Check source format is valid. It must be like domain/user/repo/path@ref
		parts := strings.Split(dep.Source, "@")
		if len(parts) != 2 {
			return ErrInvalidSource(dep.Source)
		}

		pathParts := strings.Split(parts[0], "/")
		if len(pathParts) < 3 {
			return ErrInvalidSource(dep.Source)
		}

		// Check domain is not empty
		domain := pathParts[0]
		if domain == "" {
			return ErrSourceDomainEmpty(dep.Source)
		}

		// Check if domain has been registered before or is a valid API domain if not return an error
		if _, ok := configDomains[domain]; !ok || !slices.Contains(domainAllowedApi, domain) {
			return ErrSourceUnregisteredDomain(dep.Source, domain)
		}

		// Check ref is not empty
		ref := parts[1]
		if ref == "" {
			return ErrorInvalidRef(ref, dep.Source)
		}
	}

	return nil
}
