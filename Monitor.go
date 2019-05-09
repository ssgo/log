package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Monitor(name, target, targetInfo, expect, result string, succeed bool, usedTime float32, memo string, extra ...interface{}) {
	if !logger.CheckLevel(INFO) {
		return
	}

	logger.Log(logger.MakeMonitorLog(standard.LogTypeMonitor, name, target, targetInfo, expect, result, succeed, usedTime, memo, extra...))
}

func (logger *Logger) MakeMonitorLog(logType, name, target, targetInfo, expect, result string, succeed bool, usedTime float32, memo string, extra ...interface{}) standard.MonitorLog {
	return standard.MonitorLog{
		BaseLog:    logger.MakeBaseLog(logType, extra...),
		Name:       name,
		Target:     target,
		TargetInfo: targetInfo,
		Expect:     expect,
		Result:     result,
		Succeed:    succeed,
		UsedTime:   usedTime,
		Memo:       memo,
	}
}
