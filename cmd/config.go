package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// init config by config file
func init() {
	if err := readConfig(&config, defaultConfigFilepath()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Config struct {
	BearerToken string
}

var config Config

func defaultConfigFilepath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".twego/config.toml")
}

// read config-file, and set config as it
func readConfig(config *Config, path string) error {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return createFile(path)
		}
		return err
	}

	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	return nil
}

// write config-file
func writeConfig(config Config, path string) error {
	viper.Set("BearerToken", config.BearerToken)

	return viper.WriteConfigAs(path)
}

func createFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0666); err != nil {
		return err
	}

	_, err := os.Create(path)
	return err
}
