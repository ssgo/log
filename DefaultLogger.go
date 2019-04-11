package log

var defaultLogger = Logger{}

func init() {
	defaultLogger.SetTruncations("github.com/", "/ssgo/")
}

func SetLevel(level LevelType) {
	defaultLogger.SetLevel(level)
}

func SetWriter(writer func(string)) {
	defaultLogger.SetWriter(writer)
}

func SetTruncations(truncations ...string) {
	defaultLogger.SetTruncations(truncations...)
}

func Debug(logType string, data ...interface{}) {
	defaultLogger.Debug(logType, data...)
}

func Info(logType string, data ...interface{}) {
	defaultLogger.Info(logType, data...)
}

func Warning(logType string, data ...interface{}) {
	defaultLogger.Warning(logType, data...)
}

func Error(logType string, data ...interface{}) {
	defaultLogger.Error(logType, data...)
}

func LogRequest(app, node, clientIp, fromApp, fromNode, clientId, sessionId, requestId, host string, authLevel, priority int, method, path string, requestHeaders map[string]string, requestData map[string]interface{}, usedTime float32, responseCode int, responseHeaders map[string]string, responseDataLength uint, responseData interface{}, extraInfo map[string]interface{}){
	defaultLogger.LogRequest(app, node, clientIp, fromApp, fromNode, clientId, sessionId, requestId, host, authLevel, priority, method, path, requestHeaders, requestData, usedTime, responseCode, responseHeaders, responseDataLength, responseData, extraInfo)
}
