package logger

import "log/slog"

func LogFilteredStatusCode(status int, prefix string, attrs ...any) {
	switch {
	case status >= 500:
		slog.Error(prefix+" failed", attrs...)
	case status >= 400:
		slog.Warn(prefix+" warning", attrs...)
	default:
		slog.Info(prefix+" completed", attrs...)
	}
}
