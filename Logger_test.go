package log_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ssgo/log"
	"github.com/ssgo/standard"
	"github.com/ssgo/u"
	log2 "log"
	"os"
	"strings"
	"testing"
	"time"
)


func TestLogInfoLevel(t *testing.T) {
	logFile := "tmp_test_info_level.Log"
	_ = os.Remove(logFile)
	logger := log.NewLogger(log.Config{
		File: logFile,
	})
	logger.Debug("Test", "level", 1)
	logger.Info("Test", "level", 2)
	logger.Warning("Test", "level", 3)
	logger.Error("Test", "level", 4)
	lines, _ := u.ReadFile(logFile)
	if len(lines) != 4 || strings.Index(lines[0], "info") == -1 {
		t.Error("info level test failed")
	}
	if len(lines) != 4 || strings.Index(lines[1], "warning") == -1 {
		t.Error("warning level test failed")
	}
	if len(lines) != 4 || strings.Index(lines[2], "error") == -1 {
		t.Error("error level test failed")
	}
	_ = os.Remove(logFile)
}

func TestLogDebugLevel(t *testing.T) {
	logFile := "tmp_test_debug_level.Log"
	_ = os.Remove(logFile)
	logger := log.NewLogger(log.Config{
		File:  logFile,
		Level: "debug",
	})
	logger.Debug("Test", "level", 1)
	logger.Info("Test", "level", 2)
	logger.Warning("Test", "level", 3)
	logger.Error("Test", "level", 4)
	lines, _ := u.ReadFile(logFile)
	if len(lines) != 5 || strings.Index(lines[0], "debug") == -1 {
		t.Error("info level test failed")
	}
	if len(lines) != 5 || strings.Index(lines[1], "info") == -1 {
		t.Error("info level test failed")
	}
	if len(lines) != 5 || strings.Index(lines[2], "warning") == -1 {
		t.Error("warning level test failed")
	}
	if len(lines) != 5 || strings.Index(lines[3], "error") == -1 {
		t.Error("error level test failed")
	}
	_ = os.Remove(logFile)
}

func TestLogSensitive(t *testing.T) {
	logFile := "tmp_test_sensitive.Log"
	_ = os.Remove(logFile)
	logger := log.NewLogger(log.Config{
		File:           logFile,
		Sensitive:      "phone,password,name,token,accessToken",
		RegexSensitive: "(^|[^\\d])(1\\d{10})([^\\d]|$),\\[(\\w+)\\]",
		SensitiveRule:  "12:4*4,11:3*4,7:2*2,3:1*1,2:1*0",
	})

	tests := []interface{}{
		"password", "abcd1234", "ab****34",
		"name", "张三", "张*",
		"name", "张小三", "张*三",
		"accessToken", "1122", "1**2",
		"phone", "13912345678", "139****5678",
		"phone", 13912345678, "139****5678",
		"memo", "hi, [Star]! are you ok?", "hi, [S**r]! are you ok?",
		"memo", "13912345678 is a phone, the phone is 13912345678 not 13912345677 and not 139123456781, is 13912345678", "139****5678 is a phone, the phone is 139****5678 not 139****5677 and not 139123456781, is 139****5678",
	}

	for i := 2; i < len(tests); i += 3 {
		logger.Info("Sensitive Test "+u.String(tests[i-2]), tests[i-2], tests[i-1])
	}
	lines, _ := u.ReadFile(logFile)

	lineIndex := 0
	for i := 2; i < len(tests); i += 3 {
		if strings.Index(lines[lineIndex], u.String(tests[i])) == -1 {
			t.Error("sensitive "+u.String(tests[i-2])+" test failed", lines[lineIndex], tests[i])
		}
		lineIndex ++
	}

	_ = os.Remove(logFile)
}

func TestLogSensitive2(t *testing.T) {
	logFile := "tmp_test_sensitive2.Log"
	_ = os.Remove(logFile)
	logger := log.NewLogger(log.Config{
		File:           logFile,
		Sensitive:      "phone, password, name, token, accessToken",
		RegexSensitive: "(^|[^\\d])(1\\d{10})([^\\d]|$), \\[(\\w+)\\]",
		SensitiveRule:  "12:4*4, 11:3*4, 7:2*2, 3:1*1, 2:1*0",
	})

	s := `{"app":"","authLevel":3,"clientId":"","clientIp":"127.0.0.1","fromApp":"","fromNode":"","host":"127.0.0.1:7777","logTime":1557759621.206155,"logType":"request","method":"POST","node":"192.168.199.128:7777","path":"/reg","priority":0,"proto":"1.1","requestData":{"node":"192.168.199.128:58788","token":"isBuilderTokenXXX"},"requestHeaders":{"Access-Token":"isRegisterTokenXXX","Content-Length":"77","Content-Type":"application/json","User-Agent":"Go-http-client/1.1","X-Host":"127.0.0.1:7777","X-Real-Ip":"127.0.0.1","X-Request-Id":"bqdr2uFWofNoEe6iZjJRk","X-Scheme":"http"},"requestId":"bqdr2uFWofNoEe6iZjJRk","responseCode":403,"responseData":null,"responseDataLength":0,"responseHeaders":{},"scheme":"http","serverId":"TR2TSFaHOpg","sessionId":"","traceId":"bqdr2uFWofNoEe6iZjJRk","usedTime":0.292,"type":"REJECT"}`
	m := map[string]interface{}{}
	_ = json.Unmarshal([]byte(s), &m)
	logger.Log(m)
	lines, _ := u.ReadFile(logFile)

	if len(lines) != 2 || strings.Index(lines[0], "isBu*********nXXX") == -1 || strings.Index(lines[0], "isRe**********nXXX") == -1 {
		t.Error("sensitive2 test failed", lines[0])
	}

	_ = os.Remove(logFile)
}

func TestLogRequest(t *testing.T) {
	bufw := bytes.NewBuffer([]byte{})
	log2.SetOutput(bufw)
	logger := log.New("aa112233")

	startTime := time.Now()
	time.Sleep(100 * time.Nanosecond)
	logger.Request("server1", "appA", "10.3.22.178:32421", "59.32.113.241", "appB", "10.3.22.171:12334", "HJDWAdaukhASd7", "8suAHDgsyakHU", "udaHdhagy31Dd", "abc.com", "http", "1.1", 1, 0, "POST", "/users/{userId}/events", map[string]string{"Access-Token": "abcdefg"}, map[string]interface{}{"userId": 31123}, log.MakeUesdTime(startTime, time.Now()), 200, map[string]string{"XXX": "abc"}, 3401, map[string]interface{}{"events": nil}, map[string]interface{}{"specialTag": true})
	output := bufw.String()

	//o := map[string]interface{}{}
	//_ = json.Unmarshal([]byte(output[27:]), &o)
	//o2, _ := json.MarshalIndent(o, "", "  ")
	//fmt.Println(string(o2))
	fmt.Println(output)

	if len(output) < 100 {
		t.Fatal("request len failed")
	}
	if strings.Index(output, "authLevel") == -1 {
		t.Error("request test failed")
	}
}

func TestLogMultipleInheritance(t *testing.T) {
	logFile := "tmp_test_db_error.Log"
	_ = os.Remove(logFile)
	logger := log.NewLogger(log.Config{File: logFile})
	logger.Log(standard.DBErrorLog{
		DBLog:    logger.MakeDBLog("type1", "mysql", "user:****@host:port/db", "", nil, 0),
		ErrorLog: logger.MakeErrorLog("type2", "error"),
	})

	lines, _ := u.ReadFile(logFile)
	if len(lines) != 2 || strings.Index(lines[0], "type1") == -1 {
		t.Error("multiple inheritance test failed")
	}

	_ = os.Remove(logFile)
}

func TestLogMap(t *testing.T) {
	logFile := "tmp_test_map.Log"
	_ = os.Remove(logFile)
	logger := log.NewLogger(log.Config{File: logFile})
	logger.Log(map[string]interface{}{
		"logType": "type1",
	})

	lines, _ := u.ReadFile(logFile)
	if len(lines) != 2 || strings.Index(lines[0], "type1") == -1 {
		t.Error("map test failed")
	}

	_ = os.Remove(logFile)
}
