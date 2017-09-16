package main

import (
	"os"

	"github.com/johandry/log"
	"github.com/spf13/viper"
)

func main() {
	viper.Set(log.OutputKey, os.Stderr)
	viper.Set(log.ForceColorsKey, true)
	viper.Set(log.DisableColorsKey, false)
	viper.Set(log.LevelKey, "debug")

	logMain := log.Prefix("main")
	logMain.Printf("Testing main with some parameters. %d %s %d = %d", 10, "+", 10, 10+10)

	logComp := log.Prefix("component#1")
	logComp.Debug("Debugging something here")
	logComp.Error("This is a fake error")
}
