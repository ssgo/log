package log_test

import (
	"fmt"
	"github.com/ssgo/log"
	"github.com/ssgo/standard"
	"strings"
	"testing"
	"time"
)

func TestLogLevel(t *testing.T) {
	var logger = log.Logger{}
	var logBuf = []string{}
	logger.SetWriter(func(data string) {
		logBuf = append(logBuf, data)
	})

	logger.Debug("Test", "level", 1)
	logger.Info("Test", "level", 2)
	logger.Warning("Test", "level", 3)
	logger.Error("Test", "level", 4)
	fmt.Print(logBuf)
	if len(logBuf) != 3 && strings.Index(logBuf[0], standard.LogLevelInfo) != -1 {
		t.Error("default test failed")
	}

	logBuf = []string{}
	logger.SetLevel(log.WARNING)
	logger.SetTruncations("/ssgo/")
	logger.Debug("Test", "level", 1)
	logger.Info("Test", "level", 2)
	logger.Warning("Test", "level", 3)
	logger.Error("Test", "level", 4)
	if len(logBuf) != 2 && strings.Index(logBuf[0], standard.LogLevelWarning) != -1 && strings.Index(logBuf[1], "/ssgo/") == -1 {
		t.Error("default test failed")
	}
}

func TestLogRequest(t *testing.T) {
	var logBuf = []string{}
	log.SetWriter(func(data string) {
		logBuf = append(logBuf, data)
	})

	startTime := time.Now()
	time.Sleep(100 * time.Nanosecond)
	log.LogRequest("appA", "10.3.22.178:32421", "59.32.113.241", "appB", "10.3.22.171:12334", "HJDWAdaukhASd7", "8suAHDgsyakHU", "udaHdhagy31Dd", "abc.com", 1, 2, "POST", "/users/{userId}/events", map[string]string{"Access-Token": "ab****fg"}, map[string]interface{}{"userId": 31123}, log.MakeUesdTime(startTime, time.Now()), 200, map[string]string{"XXX": "abc"}, 3401, map[string]interface{}{"events": nil}, map[string]interface{}{"specialTag": true})
	fmt.Print(logBuf)
	if len(logBuf) < 1 {
		t.Fatal("request test failed")
	}
	if len(logBuf) != 1 && strings.Index(logBuf[0], standard.LogFieldRequestAuthLevel) != -1 {
		t.Error("request test failed")
	}
}
