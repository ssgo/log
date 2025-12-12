package log

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/ssgo/standard"
	"github.com/ssgo/u"
)

func shortTime(tm string) string {
	return strings.Replace(tm[5:16], "T", " ", 1)
}

type LevelOutput struct {
	level    string
	levelKey string
}

func (levelOutput *LevelOutput) Format(v string) string {
	switch strings.ToLower(levelOutput.level) {
	case "info":
		return u.Cyan(v)
	case "warning":
		return u.Yellow(v)
	case "error":
		return u.Red(v)
	}
	return v
}

func (levelOutput *LevelOutput) BFormat(v string) string {
	switch strings.ToLower(levelOutput.level) {
	case "info":
		return u.BCyan(v)
	case "warning":
		return u.BYellow(v)
	case "error":
		return u.BRed(v)
	}
	return u.BWhite(v)
}

var errorLineMatcher = regexp.MustCompile("(\\w+\\.go:\\d+)")

func Viewable(line string) string {
	b := ParseBaseLog(line)
	if b == nil {
		// 高亮错误代码
		if strings.Contains(line, ".go:") {
			if strings.Contains(line, "/ssgo/") || strings.Contains(line, "/ssdo/") || strings.Contains(line, "/gojs/") {
				line = errorLineMatcher.ReplaceAllString(line, u.BYellow("$1"))
			} else if !strings.Contains(line, "/apigo.cc/") {
				line = errorLineMatcher.ReplaceAllString(line, u.BMagenta("$1"))
			} else if !strings.Contains(line, "/go/src/") {
				line = errorLineMatcher.ReplaceAllString(line, u.BRed("$1"))
			}
		}
		return line
	}

	var logTime time.Time
	if strings.ContainsRune(b.LogTime, 'T') {
		logTime = MakeTime(b.LogTime)
	} else {
		ft := u.Float64(b.LogTime)
		ts := int64(math.Floor(ft))
		tns := int64((ft - float64(ts)) * 1e9)
		logTime = time.Unix(ts, tns)
	}

	outs := []string{}
	t1 := strings.Split(logTime.Format("01-02 15:04:05.000"), " ")
	d := t1[0]
	t := ""
	if len(t1) > 1 {
		t = t1[1]
	}
	t2 := strings.Split(t, ".")
	s := ""
	if len(t2) > 1 {
		s = t2[1]
	}
	t = t2[0]
	outs = append(outs, u.BWhite(d+" "+t))
	if s != "" {
		outs = append(outs, u.White("."+s))
	}
	outs = append(outs, " ", u.White(b.TraceId, u.AttrDim, u.AttrUnderline))

	lo := LevelOutput{}
	if b.Extra["debug"] != nil {
		lo.level = "debug"
		lo.levelKey = "debug"
	} else if b.Extra["warning"] != nil {
		lo.level = "warning"
		lo.levelKey = "warning"
	} else if b.Extra["error"] != nil {
		lo.level = "error"
		lo.levelKey = "error"
	} else if b.Extra["info"] != nil {
		lo.level = "info"
		lo.levelKey = "info"
	} else if b.Extra["Debug"] != nil {
		lo.level = "debug"
		lo.levelKey = "Debug"
	} else if b.Extra["Warning"] != nil {
		lo.level = "warning"
		lo.levelKey = "Warning"
	} else if b.Extra["Error"] != nil {
		lo.level = "error"
		lo.levelKey = "Error"
	} else if b.Extra["Info"] != nil {
		lo.level = "info"
		lo.levelKey = "Info"
	}

	if b.LogType == standard.LogTypeRequest {
		r := standard.RequestLog{}
		ParseSpecialLog(b, &r)
		if r.ResponseCode <= 0 || (r.ResponseCode >= 400 && r.ResponseCode <= 599) {
			outs = append(outs, " ", u.BRed(u.String(r.ResponseCode)), " ", u.Red(u.String(r.UsedTime)))
		} else {
			outs = append(outs, " ", u.BGreen(u.String(r.ResponseCode)), " ", u.Green(u.String(r.UsedTime)))
		}

		outs = append(outs, "  ", r.ClientIp, u.Dim(" from"), u.Dim("("), u.Cyan(r.FromApp), u.Dim(":"), r.FromNode, u.Dim(")"), u.Dim(" to"), u.Dim("("), u.Cyan(r.App), u.Dim(":"), r.Node, u.Dim(":"), u.String(r.AuthLevel), u.Dim(":"), u.String(r.Priority), u.Dim(")"))
		if r.RequestId != r.TraceId {
			outs = append(outs, u.Dim("  requestId:"), r.RequestId)
		}
		outs = append(outs, "  ", u.Dim("user"), u.Dim(":"), u.Cyan(r.UserId), u.Dim(" sess"), u.Dim(":"), u.Cyan(r.SessionId), u.Dim(" dev"), u.Dim(":"), u.Cyan(r.DeviceId), u.Dim(" app"), u.Dim(":"), u.Cyan(r.ClientAppName), u.Dim(":"), u.Cyan(r.ClientAppVersion))
		outs = append(outs, "  ", r.Scheme, " ", r.Proto, " ", r.Host, " ", r.Method, " ", u.Cyan(r.Path))
		if r.RequestData != nil {
			for k, v := range r.RequestData {
				outs = append(outs, "  ", u.Cyan(k, u.AttrItalic), u.Dim(":"), u.String(v))
			}
		}

		if r.RequestHeaders != nil {
			for k, v := range r.RequestHeaders {
				outs = append(outs, "  ", u.Cyan(k, u.AttrDim, u.AttrItalic), u.Dim(":"), u.String(v))
			}
		}

		outs = append(outs, "  ", u.BWhite(u.String(r.ResponseDataLength)))
		if r.ResponseHeaders != nil {
			for k, v := range r.ResponseHeaders {
				outs = append(outs, "  ", u.Blue(k, u.AttrDim, u.AttrItalic), u.Dim(":"), u.String(v))
			}
		}
		outs = append(outs, "  ", u.String(r.ResponseData))
	} else if b.LogType == standard.LogTypeStatistic {
		r := standard.StatisticLog{}
		ParseSpecialLog(b, &r)
		outs = append(outs, " ", u.Cyan(r.Name, u.AttrBold))
		outs = append(outs, "  ", u.Dim(r.App))
		outs = append(outs, " ", u.Dim(shortTime(r.StartTime)+" ~ "+shortTime(r.EndTime)))
		outs = append(outs, " ", u.Green(u.String(r.Times)), " ", u.Magenta(u.String(r.Failed)))
		outs = append(outs, " ", fmt.Sprintf("%.4f", r.Min), " ", u.Cyan(fmt.Sprintf("%.4f", r.Avg)), " ", fmt.Sprintf("%.4f", r.Max))
	} else if b.LogType == standard.LogTypeTask {
		r := standard.TaskLog{}
		ParseSpecialLog(b, &r)
		if r.Succeed {
			outs = append(outs, "  ", u.Green(r.Name), " ", u.BGreen(fmt.Sprintf("%.4f", r.UsedTime)))
		} else {
			outs = append(outs, "  ", u.Red(r.Name), " ", u.BRed(fmt.Sprintf("%.4f", r.UsedTime)))
		}
		outs = append(outs, " @", u.Dim(shortTime(r.StartTime)), " @", u.Dim(r.Node))
		outs = append(outs, " ", u.Json(r.Args))
		outs = append(outs, " ", u.Magenta(r.Memo))
	} else {
		if lo.level != "" {
			outs = append(outs, " ", lo.Format(u.String(b.Extra[lo.levelKey])))
			delete(b.Extra, lo.levelKey)
		} else if b.LogType == "undefined" {
			outs = append(outs, " ", u.Dim("-"))
		} else {
			outs = append(outs, " ", u.Cyan(b.LogType, u.AttrBold))
		}
	}

	callStacks := b.Extra["callStacks"]
	if callStacks != nil {
		delete(b.Extra, "callStacks")
	}

	var codeFileMatcher = regexp.MustCompile(`(\w+?\.)(go|js)`)
	if b.Extra != nil {
		for k, v := range b.Extra {
			if k == "extra" && u.String(v)[0] == '{' {
				extra := map[string]interface{}{}
				u.UnJson(u.String(v), &extra)
				for k2, v2 := range extra {
					v2Str := u.String(v2)
					if k2 == "stack" && v2Str != "" {
						outs = append(outs, "\n")
						for _, line := range u.SplitWithoutNone(v2Str, "\n") {
							a := strings.Split(line, "```")
							for i := 0; i < len(a); i++ {
								if i%2 == 0 {
									a[i] = codeFileMatcher.ReplaceAllString(a[i], u.BRed("$1$2"))
								} else {
									a[i] = u.BCyan(a[i])
								}
							}
							outs = append(outs, "  ", strings.Join(a, ""), "\n")

						}
					} else {
						outs = append(outs, "  ", u.White(k2+":", u.AttrDim, u.AttrItalic), v2Str)
					}
				}
			} else {
				outs = append(outs, "  ", u.White(k+":", u.AttrDim, u.AttrItalic), u.String(v))
			}
		}
	}

	if callStacks != nil {
		var callStacksList []interface{}
		if callStacksStr, ok := callStacks.(string); ok && len(callStacksStr) > 2 && callStacksStr[0] == '[' {
			callStacksList = make([]interface{}, 0)
			json.Unmarshal([]byte(callStacksStr), &callStacksList)
		} else {
			callStacksList = callStacks.([]interface{})
		}

		if len(callStacksList) > 0 {
			outs = append(outs, "\n")
			for _, vi := range callStacksList {
				v := u.String(vi)
				postfix := ""
				if pos := strings.LastIndexByte(v, '/'); pos != -1 {
					postfix = v[pos+1:]
					v = v[0 : pos+1]
				} else {
					postfix = v
					v = ""
				}
				outs = append(outs, " ", u.Dim(v))
				if len(v) > 2 && (v[0] == '/' || v[1] == ':') {
					outs = append(outs, lo.BFormat(postfix))
				} else {
					outs = append(outs, lo.Format(postfix))
				}
				outs = append(outs, "\n")
			}
		} else {
			outs = append(outs, " ", lo.Format(u.String(callStacks)))
		}
	}
	return strings.Join(outs, "")
}
