package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Task(serverId, app, name string, succeed bool, usedTime float32, memo string, extra ...interface{}) {
	if !logger.CheckLevel(INFO) {
		return
	}

	logger.Log(logger.MakeTaskLog(standard.LogTypeTask, serverId, app, name, succeed, usedTime, memo, extra...))
}

func (logger *Logger) MakeTaskLog(logType, serverId, app, name string, succeed bool, usedTime float32, memo string, extra ...interface{}) standard.TaskLog {
	return standard.TaskLog{
		BaseLog:  logger.MakeBaseLog(logType, extra...),
		ServerId: serverId,
		App:      app,
		Name:     name,
		Succeed:  succeed,
		UsedTime: usedTime,
		Memo:     memo,
	}
}
