package log

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// OutputKey ...
// FilenameKey ...
// LevelKey ...
// ForceColorsKey ...
// DisableColorsKey ...
// DisableTimestampKey ...
// ShortTimestampKey ...
// TimestampFormatKey ...
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

// PrefixField ...
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

// Logger encapsulate logrus.Logger and add the prefix
type Logger struct {
	logrus.Logger

	mu     sync.Mutex
	prefix string
}

var l *Logger

func init() {
	v := viper.New()
	l = New(v)
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
	logger.SetPrefix(prefix)
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

// StdLogger return the standar logger
func StdLogger() *Logger {
	return l
}

// NewEntry create a new logger entry configured from global viper variables
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
		logLevel, err := logrus.ParseLevel(viper.GetString(LevelKey))
		if err == nil {
			logrus.SetLevel(logLevel)
		}
	}

	return logrus.NewEntry(logrus.StandardLogger())
}

// Prefix creates a new logrus.Entry from a global logrus Entry
func Prefix(prefix string) *logrus.Entry {
	return NewEntry().WithField(PrefixField, prefix)
}
