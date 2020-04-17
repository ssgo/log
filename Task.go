package log

import (
	"github.com/ssgo/standard"
	"time"
)

func (logger *Logger) Task(name, args string, succeed bool, node string, startTime time.Time, usedTime float32, memo string, extra ...interface{}) {
	if !logger.CheckLevel(INFO) {
		return
	}

	logger.Log(logger.MakeTaskLog(standard.LogTypeTask, name, args, succeed, node, startTime, usedTime, memo, extra...))
}

func (logger *Logger) MakeTaskLog(logType, name, args string, succeed bool, node string, startTime time.Time, usedTime float32, memo string, extra ...interface{}) standard.TaskLog {
	return standard.TaskLog{
		BaseLog:   logger.MakeBaseLog(logType, extra...),
		Name:      name,
		Args:      args,
		Succeed:   succeed,
		Node:      name,
		StartTime: MakeLogTime(startTime),
		UsedTime:  usedTime,
		Memo:      memo,
	}
}
