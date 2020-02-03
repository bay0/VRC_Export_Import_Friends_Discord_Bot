package config

import (
	"os"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Init reads config file for the application
func Init() {
	viper.SetConfigFile("./config.json")
	// Searches for config file in given paths and read it
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		os.Exit(3)
	}
	// Confirm which config file is used
	log.Info("Using config: " + viper.ConfigFileUsed())
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("Config file changed: ", e.Name)
	})
}