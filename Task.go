package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Task(serverId, app, name string, succeed bool, usedTime float32, memo string, extra ...interface{}) {
	if !logger.checkLevel(INFO) {
		return
	}

	logger.log(standard.TaskLog{
		BaseLog:  logger.getBaseLog(standard.LogTypeTask, extra...),
		ServerId: serverId,
		App:      app,
		Name:     name,
		Succeed:  succeed,
		UsedTime: usedTime,
		Memo:     memo,
	})
}
