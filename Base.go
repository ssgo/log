package log

import (
	"github.com/ssgo/standard"
	"github.com/ssgo/u"
	"time"
)

func (logger *Logger) Debug(debug string, extra ...interface{}) {
	if !logger.CheckLevel(DEBUG) {
		return
	}
	logger.Log(standard.DebugLog{
		BaseLog:    logger.MakeBaseLog(standard.LogTypeDebug, extra...),
		CallStacks: logger.getCallStacks(),
		Debug:      debug,
	})
}

func (logger *Logger) Info(info string, extra ...interface{}) {
	if !logger.CheckLevel(INFO) {
		return
	}
	logger.Log(standard.InfoLog{
		BaseLog: logger.MakeBaseLog(standard.LogTypeInfo, extra...),
		Info:    info,
	})
}

func (logger *Logger) Warning(warning string, extra ...interface{}) {
	if logger.CheckLevel(WARNING) {
		logger.Log(logger.MakeWarningLog(standard.LogTypeWarning, warning, extra...))
	}
}

func (logger *Logger) Error(error string, extra ...interface{}) {
	if logger.CheckLevel(ERROR) {
		logger.Log(logger.MakeErrorLog(standard.LogTypeError, error, extra...))
	}
}

func (logger *Logger) MakeDebugLog(logType, debug string, extra ...interface{}) standard.DebugLog {
	return standard.DebugLog{
		BaseLog:    logger.MakeBaseLog(logType, extra...),
		CallStacks: logger.getCallStacks(),
		Debug:      debug,
	}
}

func (logger *Logger) MakeInfoLog(logType, info string, extra ...interface{}) standard.InfoLog {
	return standard.InfoLog{
		BaseLog: logger.MakeBaseLog(logType, extra...),
		Info:    info,
	}
}

func (logger *Logger) MakeWarningLog(logType, warning string, extra ...interface{}) standard.WarningLog {
	return standard.WarningLog{
		BaseLog:    logger.MakeBaseLog(logType, extra...),
		CallStacks: logger.getCallStacks(),
		Warning:    warning,
	}
}

func (logger *Logger) MakeErrorLog(logType, error string, extra ...interface{}) standard.ErrorLog {
	return standard.ErrorLog{
		BaseLog:    logger.MakeBaseLog(logType, extra...),
		CallStacks: logger.getCallStacks(),
		Error:      error,
	}
}

func (logger *Logger) MakeBaseLog(logType string, extra ...interface{}) standard.BaseLog {
	now := time.Now()
	baseLog := standard.BaseLog{
		LogName:    logger.config.Name,
		LogTime:    MakeLogTime(now),
		LogType:    logType,
		TraceId:    logger.traceId,
		ImageName:  dockerImageName,
		ImageTag:   dockerImageTag,
		ServerName: serverName,
		ServerIp:   serverIp,
	}
	if len(extra) == 1 {
		if mapData, ok := extra[0].(map[string]string); ok {
			baseLog.Extra = mapData
			return baseLog
		}
	}
	if len(extra) > 1 {
		baseLog.Extra = map[string]string{}
		for i := 1; i < len(extra); i += 2 {
			if k, ok := extra[i-1].(string); ok {
				baseLog.Extra[k] = u.String(extra[i])
			}
		}
	}
	return baseLog
}
