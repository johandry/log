package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/Sirupsen/logrus"
	"github.com/mgutz/ansi"
)

// COLORS:
// ansi.Red
// ansi.Green
// ansi.Yellow
// ansi.Blue
// ansi.Magenta
// ansi.Cyan
// ansi.White
// ansi.LightBlack
// ansi.LightRed
// ansi.LightGreen
// ansi.LightYellow
// ansi.LightBlue
// ansi.LightMagenta
// ansi.LightCyan
// ansi.LightWhite

var (
	colorTimestamp       = ansi.LightBlack
	colorDebug           = ansi.Blue
	colorInfo            = ansi.Green
	colorWarning         = ansi.Yellow
	colorErrorFatalPanic = ansi.Red
	colorPrefix          = ansi.LightCyan
)

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

// TextFormatter ...
type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging just the time passed since beginning of execution instead of
	// the full timestamp when a TTY is attached
	ShortTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool
}

// Format ...
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	var keys = make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		if k == PrefixField {
			continue
		}
		keys = append(keys, k)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	prefixFieldClashes(entry.Data)

	isColorTerminal := checkIfTerminal(entry.Logger.Out) && (runtime.GOOS != "windows")
	isColored := (f.ForceColors || isColorTerminal) && !f.DisableColors

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}
	if isColored {
		f.printColored(b, entry, keys, timestampFormat)
	} else {
		f.printNoColor(b, entry, keys, timestampFormat)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return terminal.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func prefixFieldClashes(data logrus.Fields) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}
}

func (f *TextFormatter) levelColor(entry *logrus.Entry) (levelColor string) {
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = colorDebug
	case logrus.InfoLevel:
		levelColor = colorInfo
	case logrus.WarnLevel:
		levelColor = colorWarning
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = colorErrorFatalPanic
	}
	return
}

func (f *TextFormatter) levelText(entry *logrus.Entry) (levelText string) {
	switch entry.Level {
	case logrus.InfoLevel:
		levelText = strings.ToUpper(entry.Level.String()) + " "
	case logrus.WarnLevel:
		levelText = strings.ToUpper(entry.Level.String())[0:4] + " "
	default:
		levelText = strings.ToUpper(entry.Level.String())
	}
	return
}

func (f *TextFormatter) prefixText(entry *logrus.Entry) (prefixText string) {
	prefixText = ""
	if _, ok := entry.Data[PrefixField]; ok {
		prefixText = fmt.Sprintf(" %s:", strings.Title(entry.Data[PrefixField].(string)))
	}
	return
}

func (f *TextFormatter) timeText(entry *logrus.Entry, timestampFormat string) (timeText string) {
	timeText = entry.Time.Format(timestampFormat)
	if f.ShortTimestamp {
		timeText = fmt.Sprintf("%04d", miniTS())
	}
	return
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string) {
	levelColor := f.levelColor(entry)
	levelText := f.levelText(entry)
	prefixText := f.prefixText(entry)

	if f.DisableTimestamp {
		fmt.Fprintf(b, "%s%+5s%s%s %s", levelColor, levelText, ansi.Reset, prefixText, entry.Message)
	} else {
		timeText := f.timeText(entry, timestampFormat)
		fmt.Fprintf(b, "%s[%s]%s %s%+5s%s%s %s", colorTimestamp, timeText, ansi.Reset, levelColor, levelText, ansi.Reset, prefixText, entry.Message)
	}

	for _, k := range keys {
		v := entry.Data[k]
		f.appendKeyValue(b, fmt.Sprintf(" %s%s%s=", levelColor, k, ansi.Reset), v)
	}
}

func (f *TextFormatter) printNoColor(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat string) {
	levelText := f.levelText(entry)
	prefixText := f.prefixText(entry)

	if f.DisableTimestamp {
		fmt.Fprintf(b, "%+5s%s %s", levelText, prefixText, entry.Message)
	} else {
		timeText := f.timeText(entry, timestampFormat)
		fmt.Fprintf(b, "[%s] %+5s%s %s", timeText, levelText, prefixText, entry.Message)
	}

	for _, k := range keys {
		v := entry.Data[k]
		f.appendKeyValue(b, k, v)
	}
}

func needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	b.WriteByte(' ')
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	switch value := value.(type) {
	case string:
		if !needsQuoting(value) {
			b.WriteString(value)
		} else {
			fmt.Fprintf(b, "%q", value)
		}
	case error:
		errmsg := value.Error()
		if !needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			fmt.Fprintf(b, "%q", errmsg)
		}
	default:
		fmt.Fprint(b, value)
	}
}
