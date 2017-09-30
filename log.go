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

// Defaults values to be used when creating a Logger without user parameters
const (
	defLevel            = logrus.InfoLevel
	defForceColors      = false
	defDisableColors    = false
	defDisableTimestamp = false
	defShortTimestamp   = false
	defTimestampFormat  = ""
	defPrefix           = ""
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

// New creates a new Logger configured from an existing viper instance
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

// NewDefault creates a new Logger configured with defaults values or global
// viper values if they are defined.
func NewDefault() *Logger {
	prefix := defPrefix
	if viper.IsSet(PrefixField) {
		prefix = viper.GetString(PrefixField)
	}
	logger := &Logger{
		prefix: prefix,
	}

	forceColors := defForceColors
	if viper.IsSet(ForceColorsKey) {
		forceColors = viper.GetBool(ForceColorsKey)
	}
	disableColors := defDisableColors
	if viper.IsSet(DisableColorsKey) {
		forceColors = viper.GetBool(DisableColorsKey)
	}
	disableTimestampKey := defDisableTimestamp
	if viper.IsSet(DisableTimestampKey) {
		forceColors = viper.GetBool(DisableTimestampKey)
	}
	shortTimestamp := defShortTimestamp
	if viper.IsSet(ShortTimestampKey) {
		forceColors = viper.GetBool(ShortTimestampKey)
	}
	timestampFormat := defTimestampFormat
	if viper.IsSet(TimestampFormatKey) {
		timestampFormat = viper.GetString(TimestampFormatKey)
	}
	logger.Formatter = &TextFormatter{
		ForceColors:      forceColors,
		DisableColors:    disableColors,
		DisableTimestamp: disableTimestampKey,
		ShortTimestamp:   shortTimestamp,
		TimestampFormat:  timestampFormat,
	}
	// DisableTimestamp: true, DisableColors: true

	if viper.IsSet(OutputKey) {
		logger.Out = viper.Get(OutputKey).(io.Writer)
	} else if viper.IsSet(FilenameKey) {
		logfilename := viper.GetString(FilenameKey)
		out, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			logger.Out = out
		} else {
			logger.Errorf("Cannot create log file %s. %s", logfilename, err)
		}
	} else {
		logger.Out = os.Stderr
	}

	if viper.IsSet(LevelKey) {
		logLevel, err := logrus.ParseLevel(viper.GetString(LevelKey))
		if err == nil {
			logger.Level = logLevel
		}
	} else {
		logger.Level = defLevel
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

// Debugf redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Debugf(format, args...)
}

// Infof redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Infof(format, args...)
}

// Printf redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Printf(format, args...)
}

// Warnf redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Warnf(format, args...)
}

// Warningf redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Warnf(format, args...)
}

// Errorf redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Errorf(format, args...)
}

// Fatalf redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Fatalf(format, args...)
}

// Panicf redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Panicf(format, args...)
}

// Debug redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Debug(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Debug(args...)
}

// Print redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Print(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Info(args...)
}

// Warning redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Warning(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Warn(args...)
}

// Fatal redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Fatal(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Fatal(args...)
}

// Panic redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Panic(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Panic(args...)
}

// Debugln redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Debugln(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Debugln(args...)
}

// Infoln redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Infoln(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Infoln(args...)
}

// Println redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Println(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Println(args...)
}

// Warnln redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Warnln(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Warnln(args...)
}

// Warningln redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Warningln(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Warnln(args...)
}

// Errorln redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Errorln(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Errorln(args...)
}

// Fatalln redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Fatalln(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Fatalln(args...)
}

// Panicln redeclares the logrus method with the same name to use the prefix set
func (logger *Logger) Panicln(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Panicln(args...)
}

// Error cannot be exported as they are redefined as in mitchellh/cli.Ui to
// implement that interface
func (logger *Logger) error(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Error(args...)
}

// Info cannot be exported as they are redefined as in mitchellh/cli.Ui to
// implement that interface
func (logger *Logger) info(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Info(args...)
}

// Warn cannot be exported as they are redefined as in mitchellh/cli.Ui to
// implement that interface
func (logger *Logger) warn(args ...interface{}) {
	logger.WithField(PrefixField, logger.prefix).Warn(args...)
}
