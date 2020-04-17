package log

import (
	"encoding/json"
	"github.com/ssgo/standard"
	"github.com/ssgo/u"
	"reflect"
	"strings"
	"time"
)

func MakeTime(logTime string) time.Time {
	tm, _ := time.Parse(time.RFC3339Nano, logTime)
	return tm
}

func MakeLogTime(tm time.Time) string {
	return tm.Format(time.RFC3339Nano)
}

func MakeUesdTime(startTime, endTime time.Time) float32 {
	return float32(endTime.UnixNano()-startTime.UnixNano()) / 1e6
}

func ParseBaseLog(line string) *standard.BaseLog {
	pos := strings.IndexByte(line, '{')
	if pos == -1 {
		return ParseBadLog(line)
	} else {
		l := map[string]interface{}{}
		err := json.Unmarshal([]byte(line[pos:]), &l)
		if err != nil {
			return ParseBadLog(line)
		} else {
			baseLog := standard.BaseLog{Extra: map[string]interface{}{}}
			for k, v := range l {
				switch k {
				case "logName":
					baseLog.LogName = u.String(v)
				case "logType":
					baseLog.LogType = u.String(v)
				case "logTime":
					baseLog.LogTime = u.String(v)
				case "traceId":
					baseLog.TraceId = u.String(v)
				case "imageName":
					baseLog.ImageName = u.String(v)
				case "imageTag":
					baseLog.ImageTag = u.String(v)
				case "serverName":
					baseLog.ServerName = u.String(v)
				case "serverIp":
					baseLog.ServerIp = u.String(v)
				default:
					baseLog.Extra[k] = u.String(v)
				}
			}
			return &baseLog
		}
	}
}

func ParseBadLog(line string) *standard.BaseLog {
	baseLog := standard.BaseLog{Extra: map[string]interface{}{}}
	baseLog.LogType = standard.LogTypeUndefined
	if len(line) > 19 && line[19] == ' ' {
		tm, err := time.Parse("2006/01/02 15:04:05", line[0:19])
		if err == nil {
			baseLog.LogTime = MakeLogTime(tm)
			line = line[20:]
		} else {
			return nil
		}
	} else if len(line) > 26 && line[26] == ' ' {
		tm, err := time.Parse("2006/01/02 15:04:05.000000", line[0:26])
		if err == nil {
			baseLog.LogTime = MakeLogTime(tm)
			line = line[27:]
		} else {
			return nil
		}
	} else {
		return nil
	}
	baseLog.Extra["info"] = line
	return &baseLog
}

func ParseSpecialLog(from *standard.BaseLog, to interface{}) {
	from.Extra["logName"] = from.LogName
	from.Extra["logType"] = from.LogType
	from.Extra["logTime"] = from.LogTime
	from.Extra["traceId"] = from.TraceId
	from.Extra["imageName"] = from.ImageName
	from.Extra["imageTag"] = from.ImageTag
	from.Extra["serverName"] = from.ServerName
	from.Extra["serverIp"] = from.ServerIp
	u.Convert(from.Extra, to)
	delete(from.Extra, "logType")
	delete(from.Extra, "logTime")
	delete(from.Extra, "traceId")
	delete(from.Extra, "imageName")
	delete(from.Extra, "imageTag")
	delete(from.Extra, "serverName")
	delete(from.Extra, "serverIp")
	v := reflect.ValueOf(to)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			k := t.Field(i).Name
			delete(from.Extra, k)
			if k[0] >= 'A' && k[0] <= 'Z' {
				b := []byte(k)
				b[0] += 32
				k = string(b)
				delete(from.Extra, k)
			}
		}
	}
}
