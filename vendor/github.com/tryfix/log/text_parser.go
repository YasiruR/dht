package log

import (
	"context"
	"fmt"
	"log"
	"runtime"
)

// logParser contains parsing logic for a logger.
type logParser struct {
	*logOptions
	log *log.Logger
}

// WithPrefix appends the given prefix to the existing prefix.
func (l *logParser) WithPrefix(p string, message interface{}) string {
	if l.prefix != "" {
		if p == "" {
			return fmt.Sprintf("%s] [%+v", l.prefix, message)
		}

		return fmt.Sprintf("%s.%s] [%+v", l.prefix, p, message)
	}

	return fmt.Sprintf("%s] [%+v", p, message)
}

//isLoggable checks whether it is possible to log in the given level under current configurations.
func (l *logParser) isLoggable(level Level) bool {
	return logTypes[level] <= logTypes[l.logLevel]
}

// colored colour encodes the log level tag.
//
// Whether this returns coloured tags or not depends on the colour configuration of the logger.
func (l *logParser) colored(level Level) string {
	if l.colors {
		return string(logColors[level])
	}

	return string(level)
}

// logEntry prints the log entry to the configured io.Writer.
func (l *logParser) logEntry(ctx context.Context, level Level, message interface{}, prms ...interface{}) {
	if !l.isLoggable(level) {
		return
	}

	var params []interface{}
	format := "%s [%s] [%+v]"
	logLevel := l.colored(level)

	// add extracted trace id
	var traceID string
	if l.ctxTraceExt != nil {
		traceID = l.ctxTraceExt(ctx)
	}

	params = append(params, logLevel, traceID, fmt.Sprintf("%v", message))

	if l.filePath || l.funcPath {
		format = l.applyCallerInfo(format)
	}

	if len(prms) > 0 {
		format += " %+v"
		params = append(params, prms)
	}

	// add extracted context details
	if l.ctxExt != nil {
		if ctxData := l.ctxExt(ctx); len(ctxData) > 0 {
			format += " %v"
			params = append(params, ctxData)
		}
	}

	if level == FATAL {
		l.log.Fatalf(format, params...)
	}

	l.log.Printf(format, params...)
}

func (l *logParser) applyCallerInfo(format string) string {
	funcName := "<Unknown>"
	file := "<Unknown>"
	line := 0
	pc, f, ln, ok := runtime.Caller(l.skipFrameCount + 1)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}

	file = f
	line = ln

	// file and func format
	var filePath, funcPath string
	if l.funcPath {
		funcPath = " on func " + funcName
	}

	if l.filePath {
		filePath = " on " + file + " line " + fmt.Sprint(line)
	}

	return "%s [%s] [%+v" + funcPath + filePath + "]"
}
