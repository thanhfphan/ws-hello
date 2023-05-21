package log

import (
	"io"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var _ Logger = (*log)(nil)

type log struct {
	wrappedCores   []WrappedCore
	internalLogger *zap.Logger
}

type WrappedCore struct {
	Core           zapcore.Core
	Writer         io.WriteCloser
	WriterDisabled bool
	AtomicLevel    zap.AtomicLevel
}

func New(cfg Config) (Logger, error) {
	consoleEnc := zapcore.NewConsoleEncoder(newTermEncoderConfig(levelEncoder))
	fileEnc := zapcore.NewJSONEncoder(jsonEncoderConfig)

	consoleCore := newWrappedCore(cfg.LogLevel, os.Stdout, consoleEnc)
	rw := &lumberjack.Logger{
		Filename:   path.Join(cfg.Directory, "ws-hello.log"),
		MaxSize:    cfg.MaxSizeMb,
		MaxAge:     cfg.MaxAgeDay,
		MaxBackups: cfg.MaxFiles, // files
		Compress:   cfg.Compress,
	}
	fileCore := newWrappedCore(cfg.LogLevel, rw, fileEnc)

	l := newLogger(consoleCore, fileCore)
	return l, nil
}

func newWrappedCore(level Level, rw io.WriteCloser, encoder zapcore.Encoder) WrappedCore {
	atomicLevel := zap.NewAtomicLevelAt(zapcore.Level(level))

	core := zapcore.NewCore(encoder, zapcore.AddSync(rw), atomicLevel)
	return WrappedCore{AtomicLevel: atomicLevel, Core: core, Writer: rw}
}

func newZapLogger(wrappedCores ...WrappedCore) *zap.Logger {
	cores := make([]zapcore.Core, len(wrappedCores))
	for i, wc := range wrappedCores {
		cores[i] = wc.Core
	}
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	return logger
}

func newLogger(wrappedCores ...WrappedCore) Logger {
	return &log{
		internalLogger: newZapLogger(wrappedCores...),
		wrappedCores:   wrappedCores,
	}
}

func (l *log) Write(p []byte) (int, error) {
	for _, wc := range l.wrappedCores {
		if wc.WriterDisabled {
			continue
		}
		_, _ = wc.Writer.Write(p)
	}
	return len(p), nil
}

func (l *log) Stop() {
	for _, wc := range l.wrappedCores {
		_ = wc.Writer.Close()
	}
}

// Should only be called from [Level] functions.
func (l *log) log(level zapcore.Level, msg string, fields ...zap.Field) {
	if ce := l.internalLogger.Check(zapcore.Level(level), msg); ce != nil {
		ce.Write(fields...)
	}
}

func (l *log) Error(msg string, fields ...zap.Field) {
	l.log(zapcore.ErrorLevel, msg, fields...)
}

func (l *log) Fatal(msg string, fields ...zap.Field) {
	l.log(zapcore.FatalLevel, msg, fields...)
}

func (l *log) Warn(msg string, fields ...zap.Field) {
	l.log(zapcore.WarnLevel, msg, fields...)
}

func (l *log) Info(msg string, fields ...zap.Field) {
	l.log(zapcore.InfoLevel, msg, fields...)
}

func (l *log) Debug(msg string, fields ...zap.Field) {
	l.log(zapcore.DebugLevel, msg, fields...)
}

func (l *log) StopOnPanic() {
	if r := recover(); r != nil {
		l.Fatal("panicking", zap.Any("reason", r), zap.Stack("from"))
		l.Stop()
		panic(r)
	}
}

func (l *log) RecoverAndPanic(f func()) {
	defer l.StopOnPanic()
	f()
}

func (l *log) stopAndExit(exit func()) {
	if r := recover(); r != nil {
		l.Fatal("panicking", zap.Any("reason", r), zap.Stack("from"))
		l.Stop()
		exit()
	}
}

func (l *log) RecoverAndExit(f, exit func()) {
	defer l.stopAndExit(exit)
	f()
}
