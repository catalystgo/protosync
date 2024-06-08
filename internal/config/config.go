package config

import (
	"path"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

type (
	Dependency struct {
		Source string `yaml:"source" json:"source"`
	}

	Domain struct {
		Name string `yaml:"name" json:"name"`
		API  string `yaml:"api" json:"api"`
	}

	Config struct {
		OutDir       string       `yaml:"outDir" json:"outDir"`
		Dependencies []Dependency `yaml:"dependencies" json:"dependencies"`
		Domains      []Domain     `yaml:"domains" json:"domains"`
	}
)

var (
	cfg      Config
	loadOnce sync.Once
)

func Get() *Config {
	if isConfigEmpty() {
		return nil
	}

	return &cfg
}

func Load(configPath string, outputDir string) (_ *Config, err error) {
	loadOnce.Do(func() { err = load(configPath) })

	if err != nil {
		return nil, err
	}

	// Check if config is empty
	if isConfigEmpty() {
		return nil, ErrConfigEmpty
	}

	// If outputDir is empty, set it to the directory of the config file
	if outputDir == "" {
		outputDir = path.Dir(configPath)
	}

	// Set absolute path
	cfg.OutDir = path.Join(outputDir, cfg.OutDir)

	return &cfg, nil
}

func load(configPath string) error {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return err
	}

	cfg = config

	if err := validate(); err != nil {
		return err
	}

	return nil
}

func isConfigEmpty() bool {
	return reflect.DeepEqual(cfg, Config{})
}
