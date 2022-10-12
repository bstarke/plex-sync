package main

/*
Copyright © 2022 Brad Starkenberg brad@starkenberg.net

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
