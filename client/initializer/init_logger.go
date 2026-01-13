package initializer

import (
	"log/slog"
	"strings"

	"github.com/jjonline/go-lib-backend/logger"
	"github.com/tvb-sz/serve-swagger-ui/conf"
)

//go:noinline
func iniLogger() *logger.Logger {
	// conf.Config.Log.Level, conf.Config.Log.Path, "module"
	var (
		lvl slog.Level
		log = logger.New(&logger.Options{
			Target:    conf.Config.Log.Path,
			AddSource: false,
			UseText:   false,
			MaxSize:   0,
			MaxDays:   0,
		})
	)

	switch strings.ToLower(conf.Config.Log.Level) {
	case "error", "panic", "fatal":
		lvl = slog.LevelError
	case "warning":
		lvl = slog.LevelWarn
	case "info":
		lvl = slog.LevelInfo
	case "debug", "trace":
		lvl = slog.LevelDebug
	default:
		log.Info("unsupported log level: `" + conf.Config.Log.Level + "`, use debug instead")
		lvl = slog.LevelDebug
	}
	log.GetSlogLeveler().Set(lvl)
	return log
}
