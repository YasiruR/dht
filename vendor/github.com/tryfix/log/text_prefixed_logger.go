package log

import (
	"context"
	"fmt"
)

type prefixedLogger struct {
	logParser
}

// NewLog creates a new Logger instance using existing config from the PrefixedLogger.
//
// Default configuration values can be overridden by providing Options to the function.
func (l *prefixedLogger) NewLog(opts ...Option) Logger {
	defaults := l.logOptions.copy()
	defaults.apply(opts...)

	return &logger{
		logParser: logParser{
			logOptions: defaults,
			log:        l.log,
		},
	}
}

// NewPrefixedLog creates a new NewPrefixedLogger instance using existing config from the PrefixedLogger.
//
// Default configuration values can be overridden by providing Options to the function.
func (l *prefixedLogger) NewPrefixedLog(opts ...Option) PrefixedLogger {
	defaults := l.logOptions.copy()
	defaults.apply(opts...)

	return &prefixedLogger{
		logParser: logParser{
			logOptions: defaults,
			log:        l.log,
		},
	}
}

// ErrorContext logs with ERROR level with context and prefix.
func (l *prefixedLogger) ErrorContext(ctx context.Context, prefix string, message interface{}, params ...interface{}) {
	l.logEntry(ctx, ERROR, l.WithPrefix(prefix, message), params...)
}

// WarnContext logs with WARN level with context and prefix.
func (l *prefixedLogger) WarnContext(ctx context.Context, prefix string, message interface{}, params ...interface{}) {
	l.logEntry(ctx, WARN, l.WithPrefix(prefix, message), params...)
}

// InfoContext logs with INFO level with context and prefix.
func (l *prefixedLogger) InfoContext(ctx context.Context, prefix string, message interface{}, params ...interface{}) {
	l.logEntry(ctx, INFO, l.WithPrefix(prefix, message), params...)
}

// DebugContext logs with DEBUG level with context and prefix.
func (l *prefixedLogger) DebugContext(ctx context.Context, prefix string, message interface{}, params ...interface{}) {
	l.logEntry(ctx, DEBUG, l.WithPrefix(prefix, message), params...)
}

// TraceContext logs with TRACE level with context and prefix.
func (l *prefixedLogger) TraceContext(ctx context.Context, prefix string, message interface{}, params ...interface{}) {
	l.logEntry(ctx, TRACE, l.WithPrefix(prefix, message), params...)
}

// FatalContext logs with FATAL level with context and prefix.
func (l *prefixedLogger) FatalContext(ctx context.Context, prefix string, message interface{}, params ...interface{}) {
	l.logEntry(context.Background(), FATAL, l.WithPrefix(prefix, message), params)
}

// Error logs with ERROR level with prefix.
func (l *prefixedLogger) Error(prefix string, message interface{}, params ...interface{}) {
	l.logEntry(context.Background(), ERROR, l.WithPrefix(prefix, message), params...)
}

// Warn logs with WARN level with prefix.
func (l *prefixedLogger) Warn(prefix string, message interface{}, params ...interface{}) {
	l.logEntry(context.Background(), WARN, l.WithPrefix(prefix, message), params...)
}

// Info logs with INFO level with prefix.
func (l *prefixedLogger) Info(prefix string, message interface{}, params ...interface{}) {
	l.logEntry(context.Background(), INFO, l.WithPrefix(prefix, message), params...)
}

// Debug logs with DEBUG level with prefix.
func (l *prefixedLogger) Debug(prefix string, message interface{}, params ...interface{}) {
	l.logEntry(context.Background(), DEBUG, l.WithPrefix(prefix, message), params...)
}

// Trace logs with TRACE level with prefix.
func (l *prefixedLogger) Trace(prefix string, message interface{}, params ...interface{}) {
	l.logEntry(context.Background(), TRACE, l.WithPrefix(prefix, message), params...)
}

// Fatal logs with FATAL level with prefix.
func (l *prefixedLogger) Fatal(prefix string, message interface{}, params ...interface{}) {
	l.logEntry(context.Background(), FATAL, l.WithPrefix(prefix, message), params...)
}

// Fatalln logs with FATAL level with prefix.
func (l *prefixedLogger) Fatalln(prefix string, message interface{}, params ...interface{}) {
	l.logEntry(context.Background(), FATAL, l.WithPrefix(prefix, message), params...)
}

// Print logs with INFO level.
func (l *prefixedLogger) Print(v ...interface{}) {
	l.logEntry(context.Background(), INFO, v, l.colored(`INFO`))
}

// Printf logs with INFO level.
func (l *prefixedLogger) Printf(format string, v ...interface{}) {
	l.logEntry(context.Background(), INFO, fmt.Sprintf(format, v...), l.colored(`INFO`))
}

// Println logs with INFO level.
func (l *prefixedLogger) Println(v ...interface{}) {
	l.logEntry(context.Background(), INFO, v, l.colored(`INFO`))
}
