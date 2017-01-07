package main

import (
	"os"

	"github.com/johandry/log"
	"github.com/spf13/viper"
)

const configfilename = "config"

func initOne() {
	// Set default for the log filename
	viper.SetDefault("log_file", "default.log")
	viper.SetDefault("log_level", "info")

	// Bind the 'log' parameter to the LOG environment variable
	viper.BindEnv("log_file")
	viper.BindEnv("log_level")

	//Allow the option to set the values in a config file
	viper.SetConfigName(configfilename)
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}

func initTwo() {
	viper.Set("log_file", "set.log")
	viper.Set("log_level", "error")
}

func main() {
	// You may want to do this at the console. However, the parameter in the
	// config file has priority
	if os.Getenv("log_file") != "" {
		os.Setenv("log_file", "environment.log")
	}
	if os.Getenv("log_level") != "" {
		os.Setenv("log_level", "debug")
	}

	// initOne()
	initTwo()

	log.Prefix("main").Info("Information")
	log.Std().Debug("Debuging")
	log.Prefix("main").Warn("Warning")
	log.Prefix("main").Error("Error")
	log.Std().Fatal("Fatal")
}
