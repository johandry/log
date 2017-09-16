package main

import (
	"os"

	"github.com/johandry/log"
	"github.com/spf13/viper"
)

func main() {
	v1 := viper.New()
	v1.Set(log.OutputKey, os.Stderr)
	v1.Set(log.ForceColorsKey, true)
	v1.Set(log.DisableColorsKey, false)
	v1.Set(log.LevelKey, "debug")
	v1.Set(log.PrefixField, "main")

	logMain := log.Prefix("main")
	logMain.Printf("Testing main with some parameters. %d %s %d = %d", 10, "+", 10, 10+10)

	logComp := log.Prefix("component#1")
	logComp.Debug("Debugging something here")
	logComp.Error("This is a fake error")
}
