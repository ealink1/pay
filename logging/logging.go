package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"pay/config"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var global atomic.Pointer[zap.Logger]

func Init(cfg config.LogConfig) (*zap.Logger, error) {
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.Level))); err != nil {
		return nil, fmt.Errorf("invalid log.level: %w", err)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)
	if strings.ToLower(cfg.Encoding) == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	ws, err := buildWriteSyncer(cfg)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, ws, level)
	logger := zap.New(core)
	if cfg.Development {
		logger = logger.WithOptions(zap.Development())
	}

	global.Store(logger)
	return logger, nil
}

func L() *zap.Logger {
	if l := global.Load(); l != nil {
		return l
	}
	logger, _ := zap.NewProduction()
	global.Store(logger)
	return logger
}

func buildWriteSyncer(cfg config.LogConfig) (zapcore.WriteSyncer, error) {
	var syncers []zapcore.WriteSyncer

	if cfg.File.Enabled {
		dir := cfg.File.Dir
		if dir == "" {
			dir = "logs"
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}

		baseName := cfg.File.Filename
		if baseName == "" {
			baseName = "app"
		}

		pattern := filepath.Join(dir, baseName+"-%Y%m%d.log")
		linkName := filepath.Join(dir, baseName+".log")

		rotationTime := time.Duration(cfg.File.RotateHours) * time.Hour
		maxAge := time.Duration(cfg.File.MaxAgeDays) * 24 * time.Hour

		rl, err := rotatelogs.New(
			pattern,
			rotatelogs.WithLinkName(linkName),
			rotatelogs.WithRotationTime(rotationTime),
			rotatelogs.WithMaxAge(maxAge),
		)
		if err != nil {
			return nil, err
		}
		syncers = append(syncers, zapcore.AddSync(rl))
	}

	if len(cfg.OutputPaths) > 0 {
		if ws, _, err := zap.Open(cfg.OutputPaths...); err == nil {
			syncers = append(syncers, ws)
		} else if !cfg.File.Enabled {
			return nil, err
		}
	}

	if len(syncers) == 0 {
		ws, _, err := zap.Open("stdout")
		if err != nil {
			return nil, err
		}
		syncers = append(syncers, ws)
	}

	return zapcore.NewMultiWriteSyncer(syncers...), nil
}
