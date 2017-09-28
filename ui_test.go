package log_test

import (
	"bytes"
	"testing"

	"github.com/johandry/log"
	"github.com/mitchellh/cli"
	"github.com/spf13/viper"
)

// TestUI test all the methods defined in the interface mitchellh/cli.Ui
func TestUI(t *testing.T) {
	var expectedLogMessage string
	var actualLogMessage string

	var b bytes.Buffer
	v := viper.New()

	v.Set(log.ForceColorsKey, false)
	v.Set(log.LevelKey, "debug")
	l := newLogger(&b, v)
	l.SetPrefix("test")

	var ui cli.Ui

	ui = l

	ui.Warn("Warning")
	expectedLogMessage = " WARN  Test: Warning"
	actualLogMessage = removeTimestamp(b.String())
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()

	// Ask(string) (string, error)
	resp, err := ui.Ask("Something")
	if err != nil {
		t.Errorf("Error trying to ask something. %v", err)
	}
	expectedLogMessage = ""
	actualLogMessage = resp
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()

	// AskSecret(string) (string, error)
	resp, err = ui.AskSecret("Something")
	if err != nil {
		t.Errorf("Error trying to ask a secret. %v", err)
	}
	expectedLogMessage = ""
	actualLogMessage = resp
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()

	// Output(string)
	ui.Output("Output")
	expectedLogMessage = " INFO  Test: Output"
	actualLogMessage = removeTimestamp(b.String())
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()

	// Info(string)
	ui.Info("Info")
	expectedLogMessage = " INFO  Test: Info"
	actualLogMessage = removeTimestamp(b.String())
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()

	// Error(string)
	ui.Error("Error")
	expectedLogMessage = " ERROR Test: Error"
	actualLogMessage = removeTimestamp(b.String())
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()

	// Warn(string)
	ui.Warn("Warning")
	expectedLogMessage = " WARN  Test: Warning"
	actualLogMessage = removeTimestamp(b.String())
	if actualLogMessage != expectedLogMessage {
		t.Errorf("Expected '%s', but got '%s'", expectedLogMessage, actualLogMessage)
	}
	b.Reset()
}
