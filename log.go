package log

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Log levels:
// 	"debug"
// 	"info"
// 	"warning"
// 	"error"
// 	"fatal"
// 	"panic"

// Create a new logrus.Logger
func New(prefix string) *logrus.Entry {
	var entry *logrus.Entry

	if viper.IsSet("log_file") {
		logfilename := viper.GetString("log_file")
		out, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			logrus.SetOutput(out)
		} else {
			logrus.Errorf("Cannot create log file %s. %s", logfilename, err)
		}

	}

	if viper.IsSet("log_level") {
		log_level, err := logrus.ParseLevel(viper.GetString("log_level"))
		if err == nil {
			logrus.SetLevel(log_level)
		}
		// else {
		// 	loglevels := make([]string, len(logrus.AllLevels), cap(logrus.AllLevels))
		// 	for i, lvl := range logrus.AllLevels {
		// 		loglevels[i] = lvl.String()
		// 	}
		// 	logrus.Errorf("Unknown log_level %q. Available Log Levels are: %s", viper.GetString("log_level"), strings.Join(loglevels, ","))
		// }
	}

	if prefix != "" {
		entry = logrus.WithField("prefix", prefix)
	} else {
		entry = logrus.NewEntry(logrus.StandardLogger())
	}

	return entry
}

func Prefix(prefix string) *logrus.Entry {
	return New(prefix)
}

func Std() *logrus.Entry {
	return New("")
}
