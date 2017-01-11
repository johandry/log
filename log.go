package log

import (
	"io"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	OutputKey           = "log_output"
	FilenameKey         = "log_filename"
	LevelKey            = "log_level"
	ForceColorsKey      = "log_color"
	DisableColorsKey    = "log_nocolor"
	DisableTimestampKey = "log_notimestamp"
	ShortTimestampKey   = "log_shorttimestamp"
	TimestampFormatKey  = "log_formattimestamp"
)

const (
	PrefixField = "prefix"
)

// Log levels:
// 	"debug"  - DEBUG
// 	"info"   - INFO
// 	"warning"- WARN
// 	"error"  - ERROR
// 	"fatal"  - FATAL
// 	"panic"  - PANIC

func NewEntry() *logrus.Entry {
	logrus.SetFormatter(&TextFormatter{
		ForceColors:      viper.GetBool(ForceColorsKey),
		DisableColors:    viper.GetBool(DisableColorsKey),
		DisableTimestamp: viper.GetBool(DisableTimestampKey),
		ShortTimestamp:   viper.GetBool(ShortTimestampKey),
		TimestampFormat:  viper.GetString(TimestampFormatKey),
	})
	// DisableTimestamp: true, DisableColors: true

	if viper.IsSet(OutputKey) {
		writer := viper.Get(OutputKey).(io.Writer)
		logrus.SetOutput(writer)
	} else if viper.IsSet(FilenameKey) {
		logfilename := viper.GetString(FilenameKey)
		out, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			logrus.SetOutput(out)
		} else {
			logrus.Errorf("Cannot create log file %s. %s", logfilename, err)
		}
	}

	if viper.IsSet(LevelKey) {
		log_level, err := logrus.ParseLevel(viper.GetString(LevelKey))
		if err == nil {
			logrus.SetLevel(log_level)
		}
	}

	return logrus.NewEntry(logrus.StandardLogger())
}

// Create a new logrus.Logger
func NewEntryWithPrefix(prefix string) *logrus.Entry {
	return NewEntry().WithField(PrefixField, prefix)
}

func Prefix(prefix string) *logrus.Entry {
	return NewEntryWithPrefix(prefix)
}

func Std() *logrus.Entry {
	return NewEntry()
}
