package log

var DefaultLogger = &Logger{truncations: []string{"github.com/", "/ssgo/"}}

func New(traceId string) *Logger {
	newLogger := *DefaultLogger
	newLogger.traceId = traceId
	return &newLogger
}
