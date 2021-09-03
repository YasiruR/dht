package log

import (
	"context"
)

type jsonLogger struct {
	jsonLogParser
}

// newJsonLogger creates a new instance of json logger.
func newJsonLogger(o *logOptions) Logger {
	return &jsonLogger{
		jsonLogParser: newJsonLogParser(o),
	}
}

// NewLog creates a new Logger instance using existing config from the Logger.
//
// Default configuration values can be overridden by providing Options to the function.
func (l *jsonLogger) NewLog(opts ...Option) Logger {
	defaults := l.logOptions.copy()
	defaults.apply(opts...)

	return newJsonLogger(defaults)
}

// NewPrefixedLog creates a new NewPrefixedLogger instance using existing config from the Logger.
//
// Default configuration values can be overridden by providing Options to the function.
func (l *jsonLogger) NewPrefixedLog(opts ...Option) PrefixedLogger {
	defaults := l.logOptions.copy()
	defaults.apply(opts...)

	return newPrefixedJsonLogger(defaults)
}

// Error logs with ERROR level.
func (l *jsonLogger) Error(message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(context.Background(), l.jsonLogParser.log.Error(), "", params...).Msgf("%s", message)
}

// Warn logs with WARN level.
func (l *jsonLogger) Warn(message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(context.Background(), l.jsonLogParser.log.Warn(), "", params...).Msgf("%s", message)
}

// Info logs with INFO level.
func (l *jsonLogger) Info(message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(context.Background(), l.jsonLogParser.log.Info(), "", params...).Msgf("%s", message)
}

// Debug logs with DEBUG level.
func (l *jsonLogger) Debug(message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(context.Background(), l.jsonLogParser.log.Debug(), "", params...).Msgf("%s", message)
}

// Trace logs with TRACE level.
func (l *jsonLogger) Trace(message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(context.Background(), l.jsonLogParser.log.Trace(), "", params...).Msgf("%s", message)
}

// Fatal logs with FATAL level.
func (l *jsonLogger) Fatal(message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(context.Background(), l.jsonLogParser.log.Fatal(), "", params...).Msgf("%s", message)
}

// Fatalln logs with FATAL level.
func (l *jsonLogger) Fatalln(message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(context.Background(), l.jsonLogParser.log.Fatal(), "", params...).Msgf("%s", message)
}

// ErrorContext logs with ERROR level with context.
func (l *jsonLogger) ErrorContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(ctx, l.jsonLogParser.log.Error(), "", params...).Msgf("%s", message)
}

// WarnContext logs with WARN level with context.
func (l *jsonLogger) WarnContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(ctx, l.jsonLogParser.log.Warn(), "", params...).Msgf("%s", message)
}

// InfoContext logs with INFO level with context.
func (l *jsonLogger) InfoContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(ctx, l.jsonLogParser.log.Info(), "", params...).Msgf("%s", message)
}

// DebugContext logs with DEBUG level with context.
func (l *jsonLogger) DebugContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(ctx, l.jsonLogParser.log.Debug(), "", params...).Msgf("%s", message)
}

// TraceContext logs with TRACE level with context.
func (l *jsonLogger) TraceContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(ctx, l.jsonLogParser.log.Trace(), "", params...).Msgf("%s", message)
}

// FatalContext logs with FATAL level with context.
func (l *jsonLogger) FatalContext(ctx context.Context, message interface{}, params ...interface{}) {
	l.jsonLogParser.parse(ctx, l.jsonLogParser.log.Fatal(), "", params...).Msgf("%s", message)
}

// Print logs without a level.
func (l *jsonLogger) Print(v ...interface{}) {
	l.jsonLogParser.print(v...)
}

// Printf logs without a level.
func (l *jsonLogger) Printf(format string, v ...interface{}) {
	l.jsonLogParser.printf(format, v...)
}

// Println logs without a level.
func (l *jsonLogger) Println(v ...interface{}) {
	l.jsonLogParser.print(v...)
}
