package logger

import (
	"log/slog"
	"net/http"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func (l *Logger) LogError(r *http.Request, err error) {
	if r == nil {
		l.logger.Error(err.Error())
		return
	}
	l.logger.Error(err.Error(),
		slog.String("request_method", r.Method),
		slog.String("request_url", r.URL.String()))
}

func SetUpLogger() *Logger {
	return &Logger{slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))}
}
