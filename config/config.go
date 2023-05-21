package config

import (
	"github.com/thanhfphan/ws-hello/pkg/log"
)

type Config struct {
	LogConfig log.Config
}

func Build() (*Config, error) {
	logCfg, err := getLogConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		LogConfig: logCfg,
	}, nil
}

func getLogConfig() (log.Config, error) {
	// ex, _ := os.Executable()
	// if err != nil {
	// 	return log.Config{}, err
	// }
	// executablePath := filepath.Dir(ex)
	// logPath := filepath.Join(executablePath, "logs")
	logPath := "logs"

	cfg := log.Config{}
	cfg.Directory = logPath
	cfg.LogLevel = log.InfoLevel
	cfg.MaxSizeMb = 8
	cfg.MaxFiles = 7
	cfg.MaxAgeDay = 0 // 0 mean retain all old log
	cfg.Compress = false

	return cfg, nil
}
