package config

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/viper"
)

var (
	configFileFolder = ".config/bm/"
	// TODO: added for dev
	//	configName       = "config"
	configName = "config-dev"
	envPrefix  = "BM"

	homeDirectory  string
	configFilePath string
)

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	homeDirectory = usr.HomeDir
	configFilePath = fmt.Sprintf("%s/%s", homeDirectory, configFileFolder)
	err = createConfigFile()
	if err != nil {
		panic(err)
	}

}

func LoadConfig() error {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configFilePath)
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("could not load config file %s", err)
	}

	return nil
}

func saveDefaults() error {
	viper.SetConfigName(configName)
	viper.AddConfigPath(configFilePath)

	viper.Set("bookmarkfolder", fmt.Sprintf("%s/%s", homeDirectory, ".local/bookmarks"))
	viper.Set("addtaskcommand", "toduit create '{Title}' '{URL}' -p Inbox")
	return viper.WriteConfig()
}

func createConfigFile() error {
	filepath := fmt.Sprintf("%s%s.yaml", configFilePath, configName)
	if err := os.MkdirAll(configFilePath, 0770); err != nil {
		return err
	}

	if _, err := os.Stat(filepath); err == nil {
		return nil
	}

	if _, err := os.Create(filepath); err != nil {
		return err
	}

	return saveDefaults()
}
