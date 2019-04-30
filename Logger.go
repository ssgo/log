package log

import (
	"encoding/json"
	"fmt"
	"github.com/ssgo/standard"
	"github.com/ssgo/u"
	"log"
	"os"
	"reflect"
	"regexp"
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
	level           LevelType
	goLogger        *log.Logger
	truncations     []string
	sensitive       map[string]bool
	regexSensitive  []*regexp.Regexp
	sensitiveRule   []sensitiveRuleInfo
	desensitization func(string) string
	traceId         string
}

type Config struct {
	Level          string
	File           string
	Truncations    []string
	Sensitive      []string
	RegexSensitive []string
	SensitiveRule  []string
}

type sensitiveRuleInfo struct {
	threshold int
	leftNum   int
	rightNum  int
}

func NewLogger(conf Config) *Logger {
	logger := Logger{
		truncations: conf.Truncations,
	}

	if conf.Sensitive != nil && len(conf.Sensitive) > 0 {
		logger.sensitive = map[string]bool{}
		for _, v := range conf.Sensitive {
			logger.sensitive[fixField(v)] = true
		}
	}

	if conf.RegexSensitive != nil && len(conf.RegexSensitive) > 0 {
		logger.regexSensitive = make([]*regexp.Regexp, 0)
		for _, v := range conf.RegexSensitive {
			r, err := regexp.Compile(v)
			if err == nil {
				logger.regexSensitive = append(logger.regexSensitive, r)
			} else {
				logger.Error(err.Error())
			}
		}
		if len(logger.regexSensitive) == 0 {
			logger.regexSensitive = nil
		}
	}

	if conf.SensitiveRule != nil && len(conf.SensitiveRule) > 0 {
		logger.sensitiveRule = make([]sensitiveRuleInfo, 0)
		for _, v := range conf.SensitiveRule {
			a1 := strings.SplitN(v, ":", 2)
			if len(a1) == 2 {
				a2 := strings.SplitN(a1[1], "*", 3)
				if len(a2) == 2 {
					threshold := u.Int(a1[0])
					leftNum := u.Int(a2[0])
					rightNum := u.Int(a2[1])
					if threshold >= 0 && threshold <= 100 && leftNum >= 0 && leftNum <= 100 && rightNum >= 0 && rightNum <= 100 {
						logger.sensitiveRule = append(logger.sensitiveRule, sensitiveRuleInfo{
							threshold: threshold,
							leftNum:   leftNum,
							rightNum:  rightNum,
						})
					}
				}
			}
		}
	}

	logLevel := strings.ToLower(conf.Level)
	if logLevel == "debug" {
		logger.level = DEBUG
	} else if logLevel == "warning" {
		logger.level = WARNING
	} else if logLevel == "error" {
		logger.level = ERROR
	} else {
		logger.level = INFO
	}

	if conf.File != "" {
		fp, err := os.OpenFile(conf.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			logger.goLogger = log.New(fp, "", 0)
		} else {
			logger.Error(err.Error())
		}
	}

	return &logger
}

func (logger *Logger) SetDesensitization(f func(v string) string) {
	logger.desensitization = f
}

func (logger *Logger) New(traceId string) *Logger {
	newLogger := *logger
	newLogger.traceId = traceId
	return &newLogger
}

func (logger *Logger) checkLevel(logLevel LevelType) bool {
	settedLevel := logger.level
	if settedLevel == 0 {
		settedLevel = INFO
	}
	return logLevel >= settedLevel
}

func (logger *Logger) log(data interface{}) {
	if logger.sensitive != nil {
		logger.fixLogData("", reflect.ValueOf(data), 0)
	}
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
		if logger.goLogger == nil {
			log.Print(string(buf))
		} else {
			logger.goLogger.Print(string(buf))
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


func (logger *Logger) isSensitiveField(s string) bool {
	return logger.sensitive[fixField(s)]
}

func (logger *Logger) fixLogData(k string, v reflect.Value, level int) (*reflect.Value) {
	if level >= 10 {
		return nil
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	if t.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			newValue := logger.fixLogData(t.Field(i).Name, v.Field(i), level+1)
			if newValue != nil {
				v.Field(i).Set(*newValue)
			}
		}
	} else if t.Kind() == reflect.Map {
		for _, mk := range v.MapKeys() {
			newValue := logger.fixLogData(u.String(mk.Interface()), v.MapIndex(mk), level+1)
			if newValue != nil {
				v.SetMapIndex(mk, *newValue)
			}
		}
	} else if t.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			newValue := logger.fixLogData("", v.Index(i), level+1)
			if newValue != nil {
				v.Index(i).Set(*newValue)
			}
		}
	} else {
		if logger.isSensitiveField(k) {
			newValue := reflect.ValueOf(logger.fixValue(u.String(v.Interface())))
			return &newValue
		} else if logger.regexSensitive != nil {
			prevStr := u.String(v.Interface())
			newStr := prevStr
			for _, rx := range logger.regexSensitive {
				matchs := rx.FindAllStringSubmatchIndex(newStr, 100)
				if len(matchs) > 0 {
					for _, subMatchs := range matchs {
						if len(subMatchs) == 4 {
							m2 := logger.fixValue(newStr[subMatchs[2]:subMatchs[3]])
							newStr = fmt.Sprint(newStr[0:subMatchs[2]], m2, newStr[subMatchs[3]:])
						} else if len(subMatchs) == 8 {
							m2 := logger.fixValue(newStr[subMatchs[4]:subMatchs[5]])
							newStr = fmt.Sprint(newStr[0:subMatchs[4]], m2, newStr[subMatchs[5]:])
						}
					}
				}
			}
			if newStr != prevStr {
				newValue := reflect.ValueOf(newStr)
				return &newValue
			}
		}
	}
	return nil
}

func (logger *Logger) fixValue(s string) string {
	if logger.desensitization == nil {
		if logger.sensitiveRule != nil {
			runes := []rune(s)
			for _, sr := range logger.sensitiveRule {
				if len(runes) >= sr.threshold {
					newRunes := make([]rune, 0)
					if sr.leftNum > 0 && sr.leftNum < len(runes) {
						newRunes = append(newRunes, runes[0:sr.leftNum]...)
					}
					starNum := len(runes) - sr.leftNum - sr.rightNum
					if starNum > 0 {
						newRunes = append(newRunes, []rune(strings.Repeat("*", starNum))...)
					}
					if sr.rightNum > 0 && sr.rightNum < len(runes) {
						newRunes = append(newRunes, runes[len(runes)-sr.rightNum:]...)
					}
					return string(newRunes)
				}
			}
		}
	} else {
		return logger.desensitization(s)
	}
	return "****"
}

func fixField(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "-", ""))
}
