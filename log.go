package log

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// OutputKey is the viper variable used to define the output
// FilenameKey is the viper variable used to define log filename
// LevelKey is the viper variable used to define log level key
// ForceColorsKey is the viper variable used to define if color should be forced
// DisableColorsKey is the viper variable used to define if the colors should be
// disabled
// DisableTimestampKey is the viper variable used to define if the log timestamp
// should be disabled
// ShortTimestampKey is the viper variable used to define the log timestamp
// short format
// TimestampFormatKey is the viper variable used to define the log timestamp
// format
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

// PrefixField is the viper variable to set and get the prefix to use in the text
// formatter
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

// Logger encapsulate logrus.Logger and add the prefix. It also implements the
// interface cli.Ui from github.com/mitchellh/cli to print logs using the text
// formatter
type Logger struct {
	logrus.Logger

	mu     sync.Mutex
	prefix string
}

// New create a new Logger configured from an existing viper
func New(v *viper.Viper) *Logger {
	logger := &Logger{
		prefix: v.GetString(PrefixField),
	}
	logger.Formatter = &TextFormatter{
		ForceColors:      v.GetBool(ForceColorsKey),
		DisableColors:    v.GetBool(DisableColorsKey),
		DisableTimestamp: v.GetBool(DisableTimestampKey),
		ShortTimestamp:   v.GetBool(ShortTimestampKey),
		TimestampFormat:  v.GetString(TimestampFormatKey),
	}
	// DisableTimestamp: true, DisableColors: true

	if v.IsSet(OutputKey) {
		logger.Out = v.Get(OutputKey).(io.Writer)
	} else if v.IsSet(FilenameKey) {
		logfilename := v.GetString(FilenameKey)
		out, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			logger.Out = out
		} else {
			logger.Errorf("Cannot create log file %s. %s", logfilename, err)
		}
	}

	if v.IsSet(LevelKey) {
		logLevel, err := logrus.ParseLevel(v.GetString(LevelKey))
		if err == nil {
			logger.Level = logLevel
		}
	}

	return logger
}

// NewEntryWithPrefix creates a new logrus.Entry with a prefix.
func (logger *Logger) NewEntryWithPrefix(prefix string) *logrus.Entry {
	return logger.WithField(PrefixField, prefix)
}

// Prefix is an alias for NewEntryWithPrefix. It's used to print a message with
// a prefix
func (logger *Logger) Prefix(prefix string) *logrus.Entry {
	return logger.NewEntryWithPrefix(prefix)
}

// GetPrefix return the prefix
func (logger *Logger) GetPrefix() string {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	return logger.prefix
}

// SetPrefix sets the output prefix for the logger.
func (logger *Logger) SetPrefix(prefix string) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.prefix = prefix
}
