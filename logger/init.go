package logger

import "github.com/tryfix/log"

var Log log.Logger

func init() {
	//logLevel := log.Level(Cfg.Level)
	Log = log.Constructor.Log(log.WithColors(true), log.WithLevel(`DEBUG`), log.WithFilePath(true))
}