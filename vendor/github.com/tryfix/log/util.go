package log

import (
	"fmt"
)

// NewLog creates a new instance of the logger.
func NewLog(options ...Option) Log {
	opts := new(logOptions)
	opts.applyDefault()
	opts.apply(options...)

	switch opts.output {
	case OutJson:
		return newJsonLogImpl(opts)
	default:
		return newTextLogImpl(opts)
	}
}

// Deprecated: Do not use with the json logger.
//
// WithPrefix appends the given prefix to the existing prefix.
func WithPrefix(p string, message interface{}) string {
	return fmt.Sprintf(`%s] [%+v`, p, message)
}
