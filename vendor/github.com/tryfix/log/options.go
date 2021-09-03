package log

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Option represents a function that does one or more alterations to 'logOptions' inside the logger.
type Option func(*logOptions)

// logOptions contains all the configuration options for the logger.
type logOptions struct {
	prefix         string
	suffix         string
	colors         bool
	logLevel       Level
	filePath       bool
	funcPath       bool
	skipFrameCount int
	writer         io.Writer
	output         Output
	ctxExt         func(ctx context.Context) []interface{}
	ctxTraceExt    func(ctx context.Context) string
}

// applyDefault applies a set of predefined configurations to the logger.
func (lOpts *logOptions) applyDefault() {
	lOpts.skipFrameCount = 2
	lOpts.colors = true
	lOpts.logLevel = TRACE
	lOpts.filePath = false
	lOpts.funcPath = false
	lOpts.writer = os.Stdout
	lOpts.output = OutText
}

// copy returns a copy of existing configuration values of the logger.
func (lOpts *logOptions) copy() *logOptions {
	return &logOptions{
		prefix:         lOpts.prefix,
		suffix:         lOpts.suffix,
		colors:         lOpts.colors,
		logLevel:       lOpts.logLevel,
		filePath:       lOpts.filePath,
		funcPath:       lOpts.funcPath,
		skipFrameCount: lOpts.skipFrameCount,
		writer:         lOpts.writer,
		ctxExt:         lOpts.ctxExt,
		ctxTraceExt:    lOpts.ctxTraceExt,
	}
}

// apply applies given configuration values to the logger.
//
// This will replace existing configuration values with provided values.
func (lOpts *logOptions) apply(options ...Option) {
	for _, opt := range options {
		opt(lOpts)
	}
}

// Deprecated: use WithSkipFrameCount instead.
//
// FileDepth sets the frame count to skip when reading filepath, func path.
func FileDepth(d int) Option {
	return func(opts *logOptions) {
		opts.skipFrameCount = d
	}
}

// WithStdOut sets the log writer.
func WithStdOut(w io.Writer) Option {
	return func(opts *logOptions) {
		opts.writer = w
	}
}

// WithSkipFrameCount sets the frame count to skip when reading filepath, func path.
func WithSkipFrameCount(c int) Option {
	return func(opts *logOptions) {
		opts.skipFrameCount = c
	}
}

// WithOutput sets the output format for log entries.
func WithOutput(o Output) Option {
	return func(opts *logOptions) {
		opts.output = o
	}
}

// WithFilePath sets whether the file path is logged or not.
func WithFilePath(enabled bool) Option {
	return func(opts *logOptions) {
		opts.filePath = enabled
	}
}

// WithFilePath sets whether the file path is logged or not.
func WithFuncPath(enabled bool) Option {
	return func(opts *logOptions) {
		opts.funcPath = enabled
	}
}

// Prefixed appends the given prefix value to the existing prefix value.
func Prefixed(prefix string) Option {
	return func(opts *logOptions) {
		if opts.prefix != `` {
			opts.prefix = fmt.Sprintf(`%s.%s`, opts.prefix, prefix)
			return
		}
		opts.prefix = prefix
	}
}

// WithColors enables colours in log messages.
func WithColors(enabled bool) Option {
	return func(opts *logOptions) {
		opts.colors = enabled
	}
}

// WithLevel sets the log level.
//
// The log level is used to determine which types of logs are logged depending on the precedence of the log level.
func WithLevel(level Level) Option {
	return func(opts *logOptions) {
		opts.logLevel = level
	}
}

// WithCtxExtractor allows setting up of a function to extract values from the context.
func WithCtxExtractor(fn func(ctx context.Context) []interface{}) Option {
	return func(opts *logOptions) {
		parent := opts.ctxExt
		opts.ctxExt = func(ctx context.Context) []interface{} {
			if parent != nil {
				return append(parent(ctx), fn(ctx)...)
			}

			return fn(ctx)
		}
	}
}

// WithCtxTraceExtractor allows setting up of a function to extract trace from the context.
// Default value func(_ context.Context) string{return ""}
func WithCtxTraceExtractor(fn func(ctx context.Context) string) Option {
	return func(opts *logOptions) {
		opts.ctxTraceExt = fn
	}
}
