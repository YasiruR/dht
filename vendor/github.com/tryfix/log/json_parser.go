package log

import (
	"context"
	"fmt"
	"runtime"

	"github.com/rs/zerolog"
)

// logParser contains parsing logic for a logger.
type jsonLogParser struct {
	*logOptions
	log zerolog.Logger
}

// newJsonLogParser creates a new instance of the json parser.
func newJsonLogParser(o *logOptions) jsonLogParser {
	return jsonLogParser{
		logOptions: o,
		log:        newZerolog(o),
	}
}

// print attaches concatenated v to the message field of the json as a single string.
func (l *jsonLogParser) print(v ...interface{}) {
	l.log.Print(v...)
}

// printf attaches the format parsed string to the message field of the json.
func (l *jsonLogParser) printf(format string, v ...interface{}) {
	l.log.Printf(format, v...)
}

// parse parses all additional data.
func (l *jsonLogParser) parse(ctx context.Context, event *zerolog.Event, prefix string, params ...interface{}) *zerolog.Event {
	event = l.withPrefix(event, prefix)
	event = l.withExtractedTrace(ctx, event)
	event = l.withExtractedCtx(ctx, event)
	event = l.withParams(event, params...)
	event = l.withCallerInfo(event)

	return event
}

// withLoggerPrefix attaches the logger prefix to the event.
func (l *jsonLogParser) withPrefix(event *zerolog.Event, prefix string) *zerolog.Event {
	const key string = "prefix"

	if l.prefix != "" {
		if prefix != "" {
			return event.Str(key, l.prefix+"."+prefix)
		}

		return event.Str(key, l.prefix)
	}

	if prefix != "" {
		return event.Str(key, prefix)
	}

	return event
}

// withExtractedTrace adds the extacted trace value to the event.
func (l *jsonLogParser) withExtractedTrace(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if l.ctxTraceExt != nil {
		if trace := l.ctxTraceExt(ctx); trace != "" {
			return event.Str("trace", trace)
		}
	}

	return event
}

// withExtractedCtx adds the extacted context values to the event.
func (l *jsonLogParser) withExtractedCtx(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	if l.ctxExt != nil {
		if ctxData := l.ctxExt(ctx); len(ctxData) > 0 {
			return event.Interface("context", ctxData)
		}
	}

	return event
}

// withParams adds additional details to the event.
func (l *jsonLogParser) withParams(event *zerolog.Event, params ...interface{}) *zerolog.Event {
	if len(params) == 0 {
		return event
	}

	arr := zerolog.Arr()
	for _, param := range params {
		arr.Str(fmt.Sprintf("%+v", param))
	}

	return event.Array("params", arr)
}

// withCallerInfo adds caller info to the event.
func (l *jsonLogParser) withCallerInfo(event *zerolog.Event) *zerolog.Event {
	if !(l.funcPath || l.filePath) {
		return event
	}

	funcName := "<Unknown>"
	file := "<Unknown>"
	line := 0
	pc, f, ln, ok := runtime.Caller(l.skipFrameCount + 1)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
		file = f
		line = ln
	}

	if l.funcPath {
		event.Str("func", funcName)
	}

	if l.filePath {
		event.Str("file", file+" line "+fmt.Sprint(line))
	}

	return event
}
