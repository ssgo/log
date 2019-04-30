package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Server(info, app string, weight int, node, proto string, startTime float64, extra ...interface{}) {
	if !logger.checkLevel(INFO) {
		return
	}
	logger.log(logger.getServerLog(standard.LogTypeServer, info, app, weight, node, proto, startTime, extra...))
}

func (logger *Logger) ServerError(error, info, app string, weight int, node, proto string, startTime float64, extra ...interface{}) {
	if !logger.checkLevel(ERROR) {
		return
	}
	logger.log(standard.ServerErrorLog{
		ServerLog: logger.getServerLog(standard.LogTypeServerError, info, app, weight, node, proto, startTime, extra...),
		ErrorLog:  standard.ErrorLog{Error: error},
	})
}

func (logger *Logger) getServerLog(logType, info, app string, weight int, node, proto string, startTime float64, extra ...interface{}) standard.ServerLog {
	return standard.ServerLog{
		InfoLog:  logger.getInfoLog(logType, info, extra...),
		App:       app,
		Weight:    weight,
		Node:      node,
		Proto:     proto,
		StartTime: startTime,
	}
}
