package main

import (
	"fmt"
	"github.com/spf13/viper"
	"plex-sync/cmd"
)

func main() {
	setup()
	cmd.Execute()
}

func setup() {
	viper.SetEnvPrefix("sync")
	viper.SetConfigName("config")           // name of config file (without extension)
	viper.AddConfigPath(".")                // look for config in the working directory
	viper.AddConfigPath("$HOME/.plex-sync") // call multiple times to add many search paths
	err := viper.ReadInConfig()             // Find and read the config file
	if err != nil {                         // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	viper.AutomaticEnv()
	viper.WatchConfig()
}
