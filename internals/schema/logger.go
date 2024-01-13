package schema

import "io"

type CreateAppLoggerConfig struct {
	ConsoleWriter io.Writer
	FileWriter    io.Writer
	// KafkaWriter *logger.KafkaLogWriter
}
