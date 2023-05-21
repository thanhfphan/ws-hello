package log

import "go.uber.org/zap"

type Logger interface {
	Fatal(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)

	// Recovers a panic, logs the error, and rethrows the panic.
	StopOnPanic()
	// If a function panics, this will log that panic and then re-panic ensuring
	// that the program logs the error before exiting.
	RecoverAndPanic(f func())
	// If a function panics, this will log that panic and then call the exit
	// function, ensuring that the program logs the error, recovers, and
	// executes the desired exit function
	RecoverAndExit(f, exit func())
	// Stop this logger and write back all meta-data.
	Stop()
}
