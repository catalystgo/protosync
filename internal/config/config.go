package config

import (
	"path"
	"sync"

	"github.com/spf13/viper"
)

type (
	Dependency struct {
		Path string `yaml:"path" json:"path"`

		Source  string   `yaml:"source" json:"source"`
		Sources []string `yaml:"sources" json:"sources"`
	}

	Domain struct {
		Host string `yaml:"host" json:"host"`
		API  string `yaml:"api" json:"api"`
	}

	Config struct {
		AbsOutDir string // absolute path of the output directory

		Directory    string        `yaml:"directory" json:"directory"`
		Dependencies []*Dependency `yaml:"dependencies" json:"dependencies"`
		Domains      []*Domain     `yaml:"domains" json:"domains"`
	}
)

var (
	cfg      Config
	loadOnce sync.Once
)

func Load(configPath string, outputDir string) (_ *Config, err error) {
	loadOnce.Do(func() { err = load(configPath, outputDir) })

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func load(configPath string, outputDir string) error {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	cfg = config

	completeConfig(configPath, outputDir)

	if err := validate(); err != nil {
		return err
	}

	return nil
}

func completeConfig(configPath string, outputDir string) {
	// If outputDir is empty, set it to the directory of the config file
	if outputDir == "" {
		outputDir = path.Dir(configPath)
	}

	// Set absolute path
	cfg.AbsOutDir = path.Join(outputDir, cfg.Directory)
}
