package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/johandry/log"
	"github.com/spf13/viper"
)

const configfilename = "config"

func initOne() {
	// Set default for the log filename
	viper.SetDefault(log.FilenameKey, "default.log")
	viper.SetDefault(log.LevelKey, "info")

	// Bind the 'log' parameter to the LOG environment variable
	viper.BindEnv(log.FilenameKey)
	viper.BindEnv(log.LevelKey)

	//Allow the option to set the values in a config file
	viper.SetConfigName(configfilename)
	viper.AddConfigPath(".")
	viper.ReadInConfig()
}

func initTwo() {
	viper.Set(log.FilenameKey, "set.log")
	viper.Set(log.LevelKey, "debug")
}

func main() {
	// You may want to do this at the console. However, the parameter in the
	// config file has priority
	if os.Getenv(log.FilenameKey) != "" {
		os.Setenv(log.FilenameKey, "environment.log")
	}
	if os.Getenv(log.LevelKey) != "" {
		os.Setenv(log.LevelKey, "debug")
	}

	// initOne()
	initTwo()
	viper.Set(log.OutputKey, os.Stderr)

	log.Prefix("main").WithFields(logrus.Fields{"key": "value", "env": "test testea"}).Info("Information")
	log.Std().Debug("Debuging")
	log.Prefix("main").Warn("Warning")
	log.Prefix("main").Error("Error")
	log.Std().Fatal("Fatal")
}
