package config

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/spf13/viper"

	"github.com/evcraddock/bm/pkg/utils"
)

const (
	ConfigFileName = "config.yaml"
)

var HomeDirectory string
var ConfigFilePath string

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	HomeDirectory = usr.HomeDir
	ConfigFilePath = HomeDirectory + "/.config/bm/"
}

type Config struct {
	BookmarkFolder string `yaml:"bookmarkFolder"`
}

func ConfigFileExists() bool {
	if configfile, err := os.Stat(path.Join(ConfigFilePath, ConfigFileName)); err != nil || configfile.IsDir() {
		return false
	}

	return true
}

func LoadConfigFile() (*Config, error) {
	viper.SetConfigFile(path.Join(ConfigFilePath, ConfigFileName))
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		BookmarkFolder: viper.GetString("config.bookmarkFolder"),
	}, nil
}

func SaveConfigFile(cfg *Config) error {
	if ok := utils.CreateFile(ConfigFilePath, ConfigFileName); !ok {
		return fmt.Errorf("Unable to create config file")
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(ConfigFilePath)
	viper.Set("config.bookmarkFolder", cfg.BookmarkFolder)
	return viper.WriteConfig()
}

func CreateDefaultConfig() error {
	return SaveConfigFile(&Config{
		BookmarkFolder: HomeDirectory + "/.local/bookmarks",
	})
}
