package log

import (
	"github.com/logrusorgru/aurora"
)

type Level string
type Output string

const (
	FATAL Level = `FATAL`
	ERROR Level = `ERROR`
	WARN  Level = `WARN`
	INFO  Level = `INFO`
	DEBUG Level = `DEBUG`
	TRACE Level = `TRACE`
)

const (
	OutText Output = `text`
	OutJson Output = `json`
)

var logColors = map[Level]string{
	FATAL: aurora.BgRed(`[FATAL]`).String(),
	ERROR: aurora.BgRed(`[ERROR]`).String(),
	WARN:  aurora.BgYellow(`[ WARN]`).String(),
	INFO:  aurora.BgBlue(`[ INFO]`).String(),
	DEBUG: aurora.BgCyan(`[DEBUG]`).String(),
	TRACE: aurora.BgMagenta(`[TRACE]`).String(),
}

var logTypes = map[Level]int{
	FATAL: 0,
	ERROR: 1,
	WARN:  2,
	INFO:  3,
	DEBUG: 4,
	TRACE: 5,
}
