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
	configFileName = "config.yaml"
)

// HomeDirectory the users home directory
var HomeDirectory string

// ConfigFilePath the path to the config file
var ConfigFilePath string

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	HomeDirectory = usr.HomeDir
	ConfigFilePath = HomeDirectory + "/.config/bm/"
}

// Config configuration data
type Config struct {
	BookmarkFolder string `yaml:"bookmarkFolder"`
	AddTaskCommand string `yaml:"addTaskCommand"`
}

// FileExists returns a boolean for whether the config file exists
func FileExists() bool {
	if configfile, err := os.Stat(path.Join(ConfigFilePath, configFileName)); err != nil || configfile.IsDir() {
		return false
	}

	return true
}

// LoadConfigFile loads the config file from the ConfigFilePath
func LoadConfigFile() (*Config, error) {
	viper.SetConfigFile(path.Join(ConfigFilePath, configFileName))
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		BookmarkFolder: viper.GetString("config.bookmarkFolder"),
		AddTaskCommand: viper.GetString("config.addTaskCommand"),
	}, nil
}

// SaveConfigFile saves the config file to the ConfigFilePath
func SaveConfigFile(cfg *Config) error {
	if ok := utils.CreateFile(ConfigFilePath, configFileName); !ok {
		return fmt.Errorf("Unable to create config file")
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(ConfigFilePath)
	viper.Set("config.bookmarkFolder", cfg.BookmarkFolder)
	return viper.WriteConfig()
}

// CreateDefaultConfig creates a default config file if one does not already exist
func CreateDefaultConfig() error {
	return SaveConfigFile(&Config{
		BookmarkFolder: HomeDirectory + "/.local/bookmarks",
	})
}
