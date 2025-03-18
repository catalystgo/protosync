package config

import (
	"path"
	"strings"

	"github.com/catalystgo/protosync/internal/domain"
)

// validate validates the config
func validate() error {
	// Check if config outputDir is empty
	if cfg.Directory == "" {
		return ErrDirectoryEmpty
	}

	// configDomains is used to store domains to check if there are duplicated domains
	configDomains := make(map[string]struct{}, len(cfg.Domains)+3)
	configDomains[domain.DefaultDomainGithub] = struct{}{}
	configDomains[domain.DefaultDomainGitlab] = struct{}{}
	configDomains[domain.DefaultDomainBitbucket] = struct{}{}

	if err := validateDomain(cfg.Domains, configDomains); err != nil {
		return err
	}

	if err := validateDependencies(configDomains); err != nil {
		return err
	}

	return nil
}

// validateDomain validates the domain api, it checks:
// - Domain is having an allowed API
// - Domain is not duplicated
func validateDomain(domains []*Domain, m map[string]struct{}) error {
	for _, d := range domains {
		// validate domain Domain
		if d.Host == "" {
			return ErrDomainEmpty
		}

		// Validate domain API
		if !domain.IsDomainAPIValid(d.API) {
			return ErrInvalidDomainAPI(d.Host)
		}

		// Check if domain has been registered before, if so return an error
		if _, ok := m[d.Host]; ok {
			return ErrDomainAlreadyExists(d.Host)
		}
		m[d.Host] = struct{}{}
	}
	return nil
}

// validateDependencies validates the dependencies, it checks:
func validateDependencies(configDomains map[string]struct{}) error {
	for idx, dep := range cfg.Dependencies {
		// Check path is under the out_dir
		fileOutputPath := path.Clean(path.Join(cfg.AbsOutDir, dep.Path))
		if !strings.HasPrefix(fileOutputPath, cfg.AbsOutDir) {
			return ErrPathNotUnderOutDir(idx, dep.Path, cfg.AbsOutDir)
		}

		// Check if source and sources are set at the same time
		if dep.Source != "" && len(dep.Sources) > 0 {
			return ErrSourceAndSourcesSet(idx)
		}

		if len(dep.Sources) == 0 {
			dep.Sources = append(dep.Sources, dep.Source)
		}

		// Validate sources
		for _, source := range dep.Sources {
			s, err := domain.ParseFile(source)
			if err != nil {
				return ErrSourceInvalid(idx, source, err)
			}

			// Check if domain has been registered before, if not return an error
			if _, ok := configDomains[s.Domain]; !ok {
				return ErrSourceUnregisteredDomain(source, s.Domain)
			}
		}
	}

	return nil
}
