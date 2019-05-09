package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Server(info, app string, weight int, node, proto string, startTime float64, extra ...interface{}) {
	if !logger.CheckLevel(INFO) {
		return
	}
	logger.Log(logger.MakeServerLog(standard.LogTypeServer, info, app, weight, node, proto, startTime, extra...))
}

func (logger *Logger) ServerError(error, app string, weight int, node, proto string, startTime float64, extra ...interface{}) {
	if !logger.CheckLevel(ERROR) {
		return
	}

	logger.Log(logger.MakeServerErrorLog(standard.LogTypeServerError, error, app, weight, node, proto, startTime, extra...))
}

func (logger *Logger) MakeServerLog(logType, info, app string, weight int, node, proto string, startTime float64, extra ...interface{}) standard.ServerLog {
	return standard.ServerLog{
		InfoLog:   logger.MakeInfoLog(logType, info, extra...),
		App:       app,
		Weight:    weight,
		Node:      node,
		Proto:     proto,
		StartTime: startTime,
	}
}

func (logger *Logger) MakeServerErrorLog(logType, error, app string, weight int, node, proto string, startTime float64, extra ...interface{}) standard.ServerErrorLog {
	return standard.ServerErrorLog{
		ErrorLog:  logger.MakeErrorLog(logType, error, extra...),
		App:       app,
		Weight:    weight,
		Node:      node,
		Proto:     proto,
		StartTime: startTime,
	}
}
