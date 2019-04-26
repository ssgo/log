package log

import (
	"encoding/json"
	"fmt"
	"github.com/ssgo/standard"
	"github.com/ssgo/u"
	"log"
	"runtime"
	"strings"
	"time"
)

type LevelType int

const DEBUG LevelType = 1
const INFO LevelType = 2
const WARNING LevelType = 3
const ERROR LevelType = 4

type Logger struct {
	level       LevelType
	truncations []string
	traceId     string
	writer      func(string)
}

func (logger *Logger) New(traceId string) Logger {
	newLogger := *logger
	newLogger.traceId = traceId
	return newLogger
}

func (logger *Logger) SetLevel(level LevelType) {
	logger.level = level
}

func (logger *Logger) SetWriter(writer func(string)) {
	logger.writer = writer
}

func (logger *Logger) SetTruncations(truncations ...string) {
	logger.truncations = append(logger.truncations, truncations...)
}

func (logger *Logger) checkLevel(logLevel LevelType) bool {
	settedLevel := logger.level
	if settedLevel == 0 {
		settedLevel = INFO
	}
	return logLevel >= settedLevel
}

func (logger *Logger) log(data interface{}) {

	buf, err := json.Marshal(data)

	if err != nil {
		// 无法序列化的数据包装为 undefined
		buf, err = json.Marshal(map[string]interface{}{
			"logTime":   MakeLogTime(time.Now()),
			"logType":   standard.LogTypeUndefined,
			"traceId":   logger.traceId,
			"undefined": fmt.Sprint(data),
		})
	}

	if err == nil {
		u.FixUpperCase(buf)
		if logger.writer == nil {
			log.Print(string(buf))
		} else {
			logger.writer(string(buf))
		}
	}
}

func (logger *Logger) getCallStacks() []string {
	callStacks := make([]string, 0)
	for i := 1; i < 20; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.Contains(file, "/go/src/") {
			continue
		}
		if strings.Contains(file, "/ssgo/log") {
			continue
		}
		if logger.truncations != nil {
			for _, truncation := range logger.truncations {
				pos := strings.Index(file, truncation)
				if pos != -1 {
					file = file[pos+len(truncation):]
				}
			}
		}
		callStacks = append(callStacks, fmt.Sprintf("%s:%d", file, line))
	}
	return callStacks
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
