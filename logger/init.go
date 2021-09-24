package logger

import "github.com/tryfix/log"

var Log log.Logger

func Init() {
	Log = log.Constructor.Log(
		log.WithColors(config.ColorsEnabled),
		log.WithLevel(log.Level(config.LogLevel)),
		log.WithFilePath(config.FilePath),
	)
}