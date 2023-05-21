package log

import (
	"go.uber.org/zap/zapcore"
)

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.CapitalString())
}
