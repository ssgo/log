package log

import (
	"github.com/ssgo/standard"
	"time"
)

func (logger *Logger) Debug(debug string, extra ...interface{}) {
	if !logger.checkLevel(DEBUG) {
		return
	}
	logger.log(standard.DebugLog{
		BaseLog:    logger.getBaseLog(standard.LogTypeDebug, extra...),
		CallStacks: logger.getCallStacks(),
		Debug:      debug,
	})
}

func (logger *Logger) Info(info string, extra ...interface{}) {
	if !logger.checkLevel(INFO) {
		return
	}
	logger.log(standard.InfoLog{
		BaseLog: logger.getBaseLog(standard.LogTypeInfo, extra...),
		Info:    info,
	})
}

func (logger *Logger) Warning(warning string, extra ...interface{}) {
	if logger.checkLevel(WARNING) {
		logger.log(logger.getWarningLog(standard.LogTypeWarning, warning, extra...))
	}
}

func (logger *Logger) Error(error string, extra ...interface{}) {
	if logger.checkLevel(ERROR) {
		logger.log(logger.getErrorLog(standard.LogTypeError, error, extra...))
	}
}

func (logger *Logger) getInfoLog(logType, info string, extra ...interface{}) standard.InfoLog {
	return standard.InfoLog{
		BaseLog: logger.getBaseLog(logType, extra...),
		Info:    info,
	}
}

func (logger *Logger) getWarningLog(logType, warning string, extra ...interface{}) standard.WarningLog {
	return standard.WarningLog{
		BaseLog: logger.getBaseLog(logType, extra...),
		CallStacks: logger.getCallStacks(),
		Warning: warning,
	}
}

func (logger *Logger) getErrorLog(logType, error string, extra ...interface{}) standard.ErrorLog {
	return standard.ErrorLog{
		BaseLog: logger.getBaseLog(logType, extra...),
		CallStacks: logger.getCallStacks(),
		Error:   error,
	}
}

func (logger *Logger) getBaseLog(logType string, extra ...interface{}) standard.BaseLog {
	baseLog := standard.BaseLog{
		LogTime: MakeLogTime(time.Now()),
		LogType: logType,
		TraceId: logger.traceId,
	}
	if len(extra) == 1 {
		if mapData, ok := extra[0].(map[string]interface{}); ok {
			baseLog.Extra = mapData
			return baseLog
		}
	}
	if len(extra) > 1 {
		baseLog.Extra = map[string]interface{}{}
		for i := 1; i < len(extra); i += 2 {
			if k, ok := extra[i-1].(string); ok {
				baseLog.Extra[k] = extra[i]
			}
		}
	}
	return baseLog
}
