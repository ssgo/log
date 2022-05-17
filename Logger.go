package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ssgo/standard"
	"github.com/ssgo/u"
	"log"
	"net"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LevelType int

const DEBUG LevelType = 1
const INFO LevelType = 2
const WARNING LevelType = 3
const ERROR LevelType = 4
const CLOSE LevelType = 5

type Logger struct {
	config          Config
	level           LevelType
	goLogger        *log.Logger
	fp              *os.File
	writer          Writer
	truncations     []string
	sensitive       map[string]bool
	regexSensitive  []*regexp.Regexp
	sensitiveRule   []sensitiveRuleInfo
	desensitization func(string) string
	traceId         string
}

type Config struct {
	Name           string
	Level          string
	File           string
	Fast           bool
	SplitTag       string
	Truncations    string
	Sensitive      string
	RegexSensitive string
	SensitiveRule  string
}

type sensitiveRuleInfo struct {
	threshold int
	leftNum   int
	rightNum  int
}

var writerMakers = make(map[string]func(*Config) Writer)

var dockerImageName = ""
var dockerImageTag = ""
var serverName = ""
var serverIp = ""

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	dockerImageName = os.Getenv("DOCKER_IMAGE_NAME")
	dockerImageTag = os.Getenv("DOCKER_IMAGE_TAG")
	serverName, _ = os.Hostname()
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, a := range addrs {
			an := a.(*net.IPNet)
			// 忽略 Docker 私有网段
			if an.IP.IsGlobalUnicast() && !strings.HasPrefix(an.IP.To4().String(), "172.17.") {
				serverIp = an.IP.To4().String()
			}
		}
	}
}

func RegisterWriterMaker(name string, f func(*Config) Writer) {
	writerMakers[name] = f
}

func NewLogger(conf Config) *Logger {
	if conf.Level == "" {
		conf.Level = "info"
	}
	if conf.Truncations == "" {
		conf.Truncations = "github.com/, golang.org/, /ssgo/"
	}
	if conf.Sensitive == "" {
		conf.Sensitive = standard.LogDefaultSensitive
	}
	if conf.SensitiveRule == "" {
		conf.SensitiveRule = "12:4*4, 11:3*4, 7:2*2, 3:1*1, 2:1*0"
	}

	if conf.Name == "" {
		// 尝试读取 $DISCOVER_APP
		conf.Name = os.Getenv("DISCOVER_APP")
		if conf.Name == "" {
			conf.Name = os.Getenv("discover_app")
			if conf.Name == "" {
				// 尝试读取 $DOCKER_IMAGE_NAME
				imageName := os.Getenv("DOCKER_IMAGE_NAME")
				a := strings.Split(imageName, "/")
				imageName = a[len(a)-1]
				imageName = strings.SplitN(imageName, ":", 2)[0]
				imageName = strings.SplitN(imageName, "#", 2)[0]
				conf.Name = imageName
				if conf.Name == "" {
					// 尝试读取进程名字
					conf.Name = path.Base(os.Args[0])
				}
			}
		}
	}

	logger := Logger{
		truncations: u.SplitTrim(conf.Truncations, ","),
	}

	if len(conf.Sensitive) > 0 {
		logger.sensitive = map[string]bool{}
		ss := u.SplitTrim(conf.Sensitive, ",")
		for _, v := range ss {
			logger.sensitive[fixField(v)] = true
		}
	}

	if len(conf.RegexSensitive) > 0 {
		logger.regexSensitive = make([]*regexp.Regexp, 0)
		ss := u.SplitTrim(conf.RegexSensitive, ",")
		for _, v := range ss {
			r, err := regexp.Compile(v)
			if err == nil {
				logger.regexSensitive = append(logger.regexSensitive, r)
			} else {
				log.Println(err.Error())
			}
		}
		if len(logger.regexSensitive) == 0 {
			logger.regexSensitive = nil
		}
	}

	if len(conf.SensitiveRule) > 0 {
		logger.sensitiveRule = make([]sensitiveRuleInfo, 0)
		ss := u.SplitTrim(conf.SensitiveRule, ",")
		for _, v := range ss {
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
		if strings.Contains(conf.File, "://") {
			writerName := strings.SplitN(conf.File, "://", 2)[0]
			m := writerMakers[writerName]
			if m != nil {
				w := writerMakers[writerName](&conf)
				if w != nil {
					logger.writer = w
					writers = append(writers, w)
				}
			} else {
				logger.Error("unsupported logger writer "+writerName, "file", conf.File)
			}
		} else {
			logFile := conf.File
			if conf.SplitTag != "" {
				// 使用切割的日志文件
				logFile += "." + time.Now().Format(conf.SplitTag)
			}
			fp, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err == nil {
				logger.fp = fp
				logger.goLogger = log.New(fp, "", log.Ldate|log.Lmicroseconds)
			} else {
				log.Println(err.Error())
			}
		}
	}
	logger.config = conf
	return &logger
}

func (logger *Logger) SetLevel(level LevelType) {
	logger.level = level
}

func (logger *Logger) Split(tag string) {
	if tag == "" {
		//tag = time.Now().Format("2006-01-02T15:04:05")
	}

	logger.goLogger = nil
	err := logger.fp.Close()
	if err != nil {
		log.Println(err.Error())
	}
	err = os.Rename(logger.config.File, logger.config.File+"."+tag)
	if err != nil {
		log.Println(err.Error())
	}

	fp, err := os.OpenFile(logger.config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		logger.fp = fp
		logger.goLogger = log.New(fp, "", log.Ldate|log.Lmicroseconds)
	} else {
		log.Println(err.Error())
	}
}

func (logger *Logger) SetDesensitization(f func(v string) string) {
	logger.desensitization = f
}

func (logger *Logger) New(traceId string) *Logger {
	newLogger := *logger
	newLogger.traceId = traceId
	return &newLogger
}

func (logger *Logger) GetTraceId() string {
	return logger.traceId
}

func (logger *Logger) CheckLevel(logLevel LevelType) bool {
	settedLevel := logger.level
	if settedLevel == 0 {
		settedLevel = INFO
	}
	return logLevel >= settedLevel
}

func flat(v reflect.Value, out reflect.Value) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Name[0] > 90 {
			continue
		}
		if t.Field(i).Anonymous && t.Field(i).Type.Kind() == reflect.Struct {
			flat(v.Field(i), out)
		} else {
			out.SetMapIndex(reflect.ValueOf(t.Field(i).Name), v.Field(i))
		}
		//else if t.Field(i).Name == "Extra" && t.Field(i).Type.Kind() == reflect.Map {
		//	for _, mk := range v.Field(i).MapKeys() {
		//		out.SetMapIndex(mk, v.Field(i).MapIndex(mk))
		//	}
		//}
	}
}

var prevDay string
var changeDayLock sync.Mutex = sync.Mutex{}

func (logger *Logger) Log(data interface{}) {

	var buf []byte
	var err error
	if !logger.config.Fast {
		// 快速模式不进行扁平化、脱敏等操作
		v := reflect.ValueOf(data)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			out := reflect.ValueOf(map[string]interface{}{})
			flat(v, out)
			v = out
		}

		if logger.sensitive != nil {
			logger.fixLogData("", v, 0)
		}
		// make extra to string
		extraKey := reflect.ValueOf("Extra")
		extraValue := v.MapIndex(extraKey)
		if extraValue.IsValid() && extraValue.CanInterface() {
			v.SetMapIndex(extraKey, reflect.ValueOf(u.String(extraValue.Interface())))
		}
		buf, err = json.Marshal(v.Interface())
	} else {
		//t1 := time.Now()
		buf, err = json.Marshal(data)
		//t2 := time.Now()
		//fmt.Println("\n\n === Marshal", float32(t2.UnixNano()-t1.UnixNano())/1000000)
		//t1 = t2
	}

	if err != nil {
		// 无法序列化的数据包装为 undefined
		buf, err = json.Marshal(map[string]interface{}{
			//"logTime":   MakeLogTime(time.Now()),
			"logType":   standard.LogTypeUndefined,
			"traceId":   logger.traceId,
			"undefined": fmt.Sprint(data),
		})
	}

	if err == nil {
		if bytes.Index(buf, []byte("Header")) != -1 {
			u.FixUpperCase(buf, []string{"Header"})
		} else {
			u.FixUpperCase(buf, nil)
		}

		if logger.writer != nil {
			if writerRunning {
				logger.writer.Log(buf)
			} else {
				log.Print("writer not running")
				log.Print(string(buf))
			}
		} else if logger.goLogger == nil {
			log.Print(string(buf))
		} else {
			// 输出到文件
			if logger.config.SplitTag != "" {
				today := time.Now().Format(logger.config.SplitTag)
				if prevDay == "" {
					changeDayLock.Lock()
					if prevDay == "" {
						prevDay = today
					}
					changeDayLock.Unlock()
				}
				if today != prevDay {
					logger.goLogger.Print("start changed log file to " + today + " from " + prevDay)
					changeDayLock.Lock()
					if today != prevDay {
						prevDay = today

						// 切换日志文件
						fp, err2 := os.OpenFile(logger.config.File+"."+today, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
						if err2 == nil {
							logger.goLogger = log.New(fp, "", log.Ldate|log.Lmicroseconds)
							if logger.fp != nil {
								_ = logger.fp.Close()
							}
							logger.fp = fp
							logger.goLogger.Print("succeed changed log file to " + today + " from " + prevDay)
						} else {
							logger.goLogger.Print("failed changed log file to " + today + " from " + prevDay)
						}
					}
					changeDayLock.Unlock()
					logger.goLogger.Print("stop changed log file to " + today + " from " + prevDay)
				}
			}
			logger.goLogger.Print(string(buf))
		}
	}
}

func (logger *Logger) getCallStacks() []string {
	callStacks := make([]string, 0)
	for i := 0; i < 50; i++ {
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
					//fmt.Println(file, file[pos+len(truncation):])
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

func (logger *Logger) fixLogData(k string, v reflect.Value, level int) *reflect.Value {
	if level >= 10 {
		return nil
	}
	isPtr := false
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
		isPtr = true
	}
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}
	if !v.IsValid() {
		return nil
	}

	t := v.Type()
	if v.Kind() == reflect.Interface {
		v = v.Elem()
		if !v.IsValid() {
			return nil
		}
		t = v.Type()
	}

	if t.Kind() == reflect.Struct {
		var newValue *reflect.Value = nil
		if !v.CanSet() {
			v2 := reflect.New(v.Type()).Elem()
			newValue = &v2
			for i := 0; i < v.NumField(); i++ {
				if t.Field(i).Name[0] > 90 {
					continue
				}
				if v2.Field(i).CanSet() {
					v2.Field(i).Set(v.Field(i))
				}
			}
			v = v2
		}
		changed := false
		for i := 0; i < v.NumField(); i++ {
			if t.Field(i).Name[0] > 90 {
				continue
			}
			newValue := logger.fixLogData(t.Field(i).Name, v.Field(i), level+1)
			if newValue != nil { // && v.Field(i).CanSet()
				changed = true
				v.Field(i).Set(*newValue)
			}
		}
		if changed {
			if isPtr {
				if newValue.CanAddr() {
					pv := newValue.Addr()
					return &pv
				} else {
					pv := reflect.New(newValue.Type())
					pv.Elem().Set(*newValue)
					return &pv
				}
			}
			return newValue
		} else {
			return nil
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
			if newValue != nil && v.Index(i).CanSet() {
				v.Index(i).Set(*newValue)
			}
		}
	} else {
		if logger.isSensitiveField(k) {
			newValue := reflect.ValueOf(logger.fixValue(u.String(v.Interface())))
			if isPtr {
				if newValue.CanAddr() {
					pv := newValue.Addr()
					return &pv
				} else {
					pv := reflect.New(newValue.Type())
					pv.Elem().Set(newValue)
					return &pv
				}
			}
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
				if isPtr {
					if newValue.CanAddr() {
						pv := newValue.Addr()
						return &pv
					} else {
						pv := reflect.New(newValue.Type())
						pv.Elem().Set(newValue)
						return &pv
					}
				}
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
