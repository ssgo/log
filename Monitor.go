package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Monitor(name, target, targetInfo, expect, result string, succeed bool, usedTime float32, memo string, extra ...interface{}) {
	if !logger.checkLevel(INFO) {
		return
	}

	logger.log(standard.MonitorLog{
		BaseLog:    logger.getBaseLog(standard.LogTypeMonitor, extra...),
		Name:       name,
		Target:     target,
		TargetInfo: targetInfo,
		Expect:     expect,
		Result:     result,
		Succeed:    succeed,
		UsedTime:   usedTime,
		Memo:       memo,
	})
}
