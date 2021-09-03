package log

import (
	"github.com/rs/zerolog"
)

type jsonLogImpl struct {
	*logOptions
	log zerolog.Logger
}

// newJsonLogImpl creates a new instance of json log implementation.
func newJsonLogImpl(opts *logOptions) Log {
	return &jsonLogImpl{
		logOptions: opts,
		log:        newZerolog(opts),
	}
}

// Log creates a new logger by extending the json logger implementation.
func (l *jsonLogImpl) Log(options ...Option) Logger {
	opts := l.logOptions.copy()
	opts.apply(options...)

	return newJsonLogger(opts)
}

// SimpleLog creates a new simple logger by extending the json logger implementation.
// *Note: not implemented, will panic.
func (*jsonLogImpl) SimpleLog() SimpleLogger {
	panic(`implement me`)
}

// PrefixedLog creates a new prefixed logger by extending the text logger implementation.
func (l *jsonLogImpl) PrefixedLog(options ...Option) PrefixedLogger {
	opts := l.logOptions.copy()
	opts.apply(options...)

	return newPrefixedJsonLogger(opts)
}

// zerologLevel convers the config log level to corresponding zerolog log level.
func zerologLevel(lvl Level) zerolog.Level {
	switch lvl {
	case FATAL:
		return zerolog.FatalLevel
	case ERROR:
		return zerolog.ErrorLevel
	case WARN:
		return zerolog.WarnLevel
	case INFO:
		return zerolog.InfoLevel
	case DEBUG:
		return zerolog.DebugLevel
	case TRACE:
		return zerolog.TraceLevel
	default:
		return zerolog.TraceLevel
	}
}

// newZerolog creates a new zerolog instance using given configs.
func newZerolog(opts *logOptions) zerolog.Logger {
	// set time format
	zerolog.TimeFieldFormat = "2006/01/02 15:04:05.000000"

	return zerolog.New(opts.writer).Level(zerologLevel(opts.logLevel)).
		With().Timestamp().Logger()
}
