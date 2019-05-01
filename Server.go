package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Server(info, app string, weight int, node, proto string, startTime float64, extra ...interface{}) {
	if !logger.checkLevel(INFO) {
		return
	}
	logger.log(standard.ServerLog{
		InfoLog:  logger.getInfoLog(standard.LogTypeServer, info, extra...),
		App:       app,
		Weight:    weight,
		Node:      node,
		Proto:     proto,
		StartTime: startTime,
	})
}

func (logger *Logger) ServerError(error, app string, weight int, node, proto string, startTime float64, extra ...interface{}) {
	if !logger.checkLevel(ERROR) {
		return
	}

	logger.log(standard.ServerErrorLog{
		ErrorLog:  logger.getErrorLog(standard.LogTypeServerError, error, extra...),
		App:       app,
		Weight:    weight,
		Node:      node,
		Proto:     proto,
		StartTime: startTime,
	})
}
