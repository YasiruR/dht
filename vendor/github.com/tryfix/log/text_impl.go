package log

import "log"

type textLogImpl struct {
	*logOptions
	log *log.Logger
}

// newTextLogImpl creates a new instance of text log implementation.
func newTextLogImpl(opts *logOptions) Log {
	return &textLogImpl{
		logOptions: opts,
		log:        log.New(opts.writer, ``, log.LstdFlags|log.Lmicroseconds),
	}
}

// Log creates a new logger by extending the text logger implementation.
func (l *textLogImpl) Log(options ...Option) Logger {
	opts := l.logOptions.copy()
	opts.apply(options...)

	return &logger{
		logParser: logParser{
			logOptions: opts,
			log:        l.log,
		},
	}
}

// SimpleLog creates a new simple logger by extending the text logger implementation.
// *Note: not implemented, will panic.
func (*textLogImpl) SimpleLog() SimpleLogger {
	panic(`implement me`)
}

// PrefixedLog creates a new prefixed logger by extending the text logger implementation.
func (l *textLogImpl) PrefixedLog(options ...Option) PrefixedLogger {
	opts := l.logOptions.copy()
	opts.apply(options...)

	return &prefixedLogger{
		logParser: logParser{
			logOptions: opts,
			log:        l.log,
		},
	}
}
