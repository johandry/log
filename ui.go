package log

// Logger also implements the interface cli.Ui from github.com/mitchellh/cli to
// print logs using the text formatter

// Ask asks the user for input using the given query. This UI do not ask
// questions so it implements nothing.
func (logger *Logger) Ask(query string) (string, error) {
	return "", nil
}

// AskSecret asks the user for input using the given query, but does not echo
// the keystrokes to the terminal. This UI do not ask questions so it implements
// nothing.
func (logger *Logger) AskSecret(query string) (string, error) {
	return "", nil
}

// Error is used for any error messages that might appear on standard error.
func (logger *Logger) Error(message string) {
	logger.Error(message)
}

// Info is called for information output.
func (logger *Logger) Info(message string) {
	logger.Info(message)
}

// Output is called for normal standard output.
func (logger *Logger) Output(message string) {
	logger.Print(message)
}

// Warn is used for any warning messages that might appear on standard error.
func (logger *Logger) Warn(message string) {
	logger.Warn(message)
}
