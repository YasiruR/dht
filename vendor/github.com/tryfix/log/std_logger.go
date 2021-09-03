package log

import "context"

var Constructor = NewLog(FileDepth(2))
var StdLogger = Constructor.Log(FileDepth(3))
var PrefixedStdLogger = Constructor.PrefixedLog(FileDepth(3))

// Fatal logs with FATAL level using the standard logger.
func Fatal(message interface{}, params ...interface{}) {
	StdLogger.Fatal(message, params...)
}

// Error logs with ERROR level using the standard logger.
func Error(message interface{}, params ...interface{}) {
	StdLogger.Error(message, params...)
}

// Warn logs with WARN level using the standard logger.
func Warn(message interface{}, params ...interface{}) {
	StdLogger.Warn(message, params...)
}

// Debug logs with DEBUG level using the standard logger.
func Debug(message interface{}, params ...interface{}) {
	StdLogger.Debug(message, params...)
}

// Info logs with INFO level using the standard logger.
func Info(message interface{}, params ...interface{}) {
	StdLogger.Info(message, params...)
}

// Trace logs with TRACE level using the standard logger.
func Trace(message interface{}, params ...interface{}) {
	StdLogger.Trace(message, params...)
}

// FatalContext logs with FATAL level with context and prefix using the standard logger.
func FatalContext(ctx context.Context, message interface{}, params ...interface{}) {
	StdLogger.FatalContext(ctx, message, params...)
}

// ErrorContext logs with ERROR level with context and prefix using the standard logger.
func ErrorContext(ctx context.Context, message interface{}, params ...interface{}) {
	StdLogger.ErrorContext(ctx, message, params...)
}

// WarnContext logs with WARN level with context and prefix using the standard logger.
func WarnContext(ctx context.Context, message interface{}, params ...interface{}) {
	StdLogger.WarnContext(ctx, message, params...)
}

// DebugContext logs with DEBUG level with context and prefix using the standard logger.
func DebugContext(ctx context.Context, message interface{}, params ...interface{}) {
	StdLogger.DebugContext(ctx, message, params...)
}

// InfoContext logs with INFO level with context and prefix using the standard logger.
func InfoContext(ctx context.Context, message interface{}, params ...interface{}) {
	StdLogger.InfoContext(ctx, message, params...)
}

// TraceContext logs with TRACE level with context and prefix using the standard logger.
func TraceContext(ctx context.Context, message interface{}, params ...interface{}) {
	StdLogger.TraceContext(ctx, message, params...)
}
