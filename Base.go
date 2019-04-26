package log

import "github.com/ssgo/standard"

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
	if !logger.checkLevel(WARNING) {
		return
	}
	logger.log(standard.WarningLog{
		BaseLog:   logger.getBaseLog(standard.LogTypeWarning, extra...),
		Warning:   warning,
	})
}

func (logger *Logger) Error(error string, extra ...interface{}) {
	if !logger.checkLevel(ERROR) {
		return
	}
	logger.log(standard.ErrorLog{
		BaseLog:    logger.getBaseLog(standard.LogTypeError, extra...),
		CallStacks: logger.getCallStacks(),
		Error:      error,
	})
}
