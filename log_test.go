package log_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/johandry/log"
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// removeTimestamp remove the timestamp form the log entry and the trailling
// newline character.
func removeTimestamp(logMessage string) string {
	noNewLine := strings.TrimSuffix(logMessage, "\n")
	noTimeStamp := strings.Split(noNewLine, "]")[1]
	noResetColor := strings.TrimPrefix(noTimeStamp, ansi.Reset)

	return noResetColor
}

func newLogger(b *bytes.Buffer, v *viper.Viper) *log.Logger {
	v.Set(log.OutputKey, b)
	l := log.New(v)
	return l
}

func TestConfigWithViperNoColor(t *testing.T) {
	var expectedLogMessage string
	var actualLogMessage string

	var b bytes.Buffer
	v := viper.New()

	v.Set(log.ForceColorsKey, false)
	v.Set(log.LevelKey, "debug")
	l := newLogger(&b, v)

	l.Prefix("test").WithFields(logrus.Fields{"key": "value", "env": "test testing"}).Info("Information")
	expectedLogMessage = " INFO  Test: Information env=\"test testing\" key=value"
	actualLogMessage = removeTimestamp(b.String())
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()

	l.Prefix("test").Warn("Warning")
	expectedLogMessage = " WARN  Test: Warning"
	actualLogMessage = removeTimestamp(b.String())
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()

	l.Prefix("test").Error("Error")
	expectedLogMessage = " ERROR Test: Error"
	actualLogMessage = removeTimestamp(b.String())
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()
}

func TestConfigWithViperColor(t *testing.T) {
	var b bytes.Buffer
	var expectedLogMessage string
	var actualLogMessage string

	v := viper.New()
	v.Set(log.OutputKey, &b)
	v.Set(log.ForceColorsKey, true)
	v.Set(log.LevelKey, "debug")

	l := log.New(v)

	// l.Prefix("test").WithFields(logrus.Fields{"key": "value", "env": "test testing"}).Info("Information")
	// expectedLogMessage = fmt.Sprintf(" %sINFO%s  Test: Information %senv%s=\"test testing\" %skey%s=value", ansi.Green, ansi.Reset, ansi.Green, ansi.Reset, ansi.Green, ansi.Reset)
	// actualLogMessage = removeTimestamp(b.String())
	// if actualLogMessage != expectedLogMessage {
	// 	t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	// }
	// b.Reset()
	//
	// l.Prefix("test").Warn("Warning")
	// expectedLogMessage = fmt.Sprintf(" %sWARN%s  Test: Warning", ansi.Yellow, ansi.Reset)
	// actualLogMessage = removeTimestamp(b.String())
	// if actualLogMessage != expectedLogMessage {
	// 	t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	// }
	// b.Reset()

	l.Prefix("test").Error("Error")
	expectedLogMessage = fmt.Sprintf(" %sERROR%s Test: Error", ansi.Red, ansi.Reset)
	actualLogMessage = removeTimestamp(b.String())
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()
}
