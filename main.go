package main

import (
	"os"

	"github.com/deepset-ai/prompthub/api"
	"github.com/deepset-ai/prompthub/index"
	"github.com/deepset-ai/prompthub/output"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// Define and parse command args
	verbosity := pflag.IntP("verbose", "v", 1, "set verbosity level: 0 silent, 1 normal, 2 debug")
	configPath := pflag.StringP("config", "c", "", "path to config file")
	help := pflag.BoolP("help", "h", false, "print args help")
	pflag.Parse()

	// Print the help message and exit if --help is passed
	if *help {
		pflag.PrintDefaults()
		os.Exit(0)
	}

	// Configure cmdline output facilities
	output.Init(*verbosity)

	// Bootstrap config, this has to be called first
	initConfig(configPath)

	// Initialize the index by reading all the prompts from file
	if err := index.Init(viper.GetString("prompts_path")); err != nil {
		os.Exit(1)
	}

	// Start the HTTP server, block until shutdown
	api.Serve()
	os.Exit(0)
}

func initConfig(configPath *string) {
	// Defaults
	viper.SetDefault("port", "80")
	viper.SetDefault("prompts_path", "./prompts")
	viper.SetDefault("allowed_origins", []string{"https://prompthub.deepset.ai"})
	// Add this line to set a default value for github_token
	viper.SetDefault("github_token", "")

	// Automatically bind all the config options to env vars
	viper.AutomaticEnv()

	if *configPath != "" {
		// Use SetConfigFile for full file paths
		viper.SetConfigFile(*configPath)
	} else {
		// Default config file lookup
		viper.SetConfigName("prompthub.yaml")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		output.INFO.Println("Configuration file not found, running with default parameters")
	} else {
		output.DEBUG.Println("Config file found at", viper.ConfigFileUsed())
	}

	// Add this block to check if the github_token is set
	if viper.GetString("github_token") == "" {
		output.INFO.Println("Warning: github_token is not set")
	}
}
