package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// NewZeroLogConsoleWriter returns new instance zerolog console logger
func NewZeroLogConsoleWriter() io.Writer {
	return zerolog.ConsoleWriter{Out: os.Stdout}
}
