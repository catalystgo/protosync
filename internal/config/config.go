package config

type Dependency struct {
	Source string `yaml:"source" json:"source"`
	Path   string `yaml:"path" json:"path"`
}

type Config struct {
	OutDir       string       `yaml:"outDir" json:"outDir"`
	Dependencies []Dependency `yaml:"dependencies" json:"dependencies"`
}
