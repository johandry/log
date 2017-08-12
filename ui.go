package log

// UI implements the interface cli.Ui from github.com/mitchellh/cli to print
// logs using the text formatter
// Prefix is the prefix that will be at the beginning of every message printed.
type UI struct {
	Prefix string
}

// Ask asks the user for input using the given query. This UI do not ask
// questions so it implements nothing.
func (u *UI) Ask(query string) (string, error) {
	return "", nil
}

// AskSecret asks the user for input using the given query, but does not echo
// the keystrokes to the terminal. This UI do not ask questions so it implements
// nothing.
func (u *UI) AskSecret(query string) (string, error) {
	return "", nil
}

// Error is used for any error messages that might appear on standard error.
func (u *UI) Error(message string) {
	Prefix(u.Prefix).Error(message)
}

// Info is called for information output.
func (u *UI) Info(message string) {
	Prefix(u.Prefix).Info(message)
}

// Output is called for normal standard output.
func (u *UI) Output(message string) {
	Prefix(u.Prefix).Print(message)
}

// Warn is used for any warning messages that might appear on standard error.
func (u *UI) Warn(message string) {
	Prefix(u.Prefix).Warn(message)
}
