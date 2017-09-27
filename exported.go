package log

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var l *Logger

func init() {
	v := viper.New()
	v.Set(ForceColorsKey, viper.GetBool(ForceColorsKey))
	v.Set(DisableColorsKey, viper.GetBool(DisableColorsKey))
	v.Set(DisableTimestampKey, viper.GetBool(DisableTimestampKey))
	v.Set(ShortTimestampKey, viper.GetBool(ShortTimestampKey))
	v.Set(TimestampFormatKey, viper.GetBool(TimestampFormatKey))

	if viper.IsSet(OutputKey) {
		v.Set(OutputKey, viper.Get(OutputKey))
	}
	if viper.IsSet(FilenameKey) {
		v.Set(FilenameKey, viper.GetString(FilenameKey))
	}
	if viper.IsSet(LevelKey) {
		v.Set(LevelKey, viper.GetString(LevelKey))
	}

	l = New(v)
}

// StdLogger return the standar logger
func StdLogger() *Logger {
	return l
}

// NewEntryWithPrefix creates a new logrus.Entry with a prefix.
func NewEntryWithPrefix(prefix string) *logrus.Entry {
	return l.NewEntryWithPrefix(prefix)
}

// Prefix is an alias for NewEntryWithPrefix. It's used to print a message with
// a different prefix
func Prefix(prefix string) *logrus.Entry {
	return l.Prefix(prefix)
}

// GetPrefix return the prefix
func GetPrefix() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.prefix
}

// SetPrefix sets the output prefix for the standard logger
func SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}
