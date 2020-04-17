package log

import (
	"github.com/ssgo/config"
)

var DefaultLogger *Logger

func init() {
	RegisterWriterMaker("es", esWriterMaker)
	RegisterWriterMaker("ess", esWriterMaker)

	conf := Config{}
	config.LoadConfig("log", &conf)
	DefaultLogger = NewLogger(conf)
}

func New(traceId string) *Logger {
	newLogger := *DefaultLogger
	newLogger.traceId = traceId
	return &newLogger
}
