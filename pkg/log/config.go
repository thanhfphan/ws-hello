package log

import "go.uber.org/zap/zapcore"

var (
	defaultEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	jsonEncoderConfig zapcore.EncoderConfig
)

func init() {
	jsonEncoderConfig = defaultEncoderConfig
	jsonEncoderConfig.EncodeLevel = levelEncoder
	jsonEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoderConfig.EncodeDuration = zapcore.NanosDurationEncoder
}

func newTermEncoderConfig(lvlEncoder zapcore.LevelEncoder) zapcore.EncoderConfig {
	config := defaultEncoderConfig
	config.EncodeLevel = lvlEncoder
	config.ConsoleSeparator = " "
	return config
}

type RotatingWriterConfig struct {
	MaxSizeMb int    `json:"max_size_mb"`
	MaxFiles  int    `json:"max_files"`
	MaxAgeDay int    `json:"max_age_day"`
	Directory string `json:"directory"`
	Compress  bool   `json:"compress"`
}

// Config defines the configuration of a logger
type Config struct {
	RotatingWriterConfig
	LogLevel Level `json:"logLevel"`
}
