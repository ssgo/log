package log_test

import (
	"fmt"
	"github.com/ssgo/log"
	"github.com/ssgo/u"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSplit(t *testing.T) {
	logFile := "split.log"
	//splitTag := "15:04:05"
	splitTag := "15:04"

	_ = os.RemoveAll(logFile + "*")

	logger := log.NewLogger(log.Config{
		File:     logFile,
		SplitTag: splitTag,
	})

	outs := make([]string, 0)
	prevTime := time.Now().Format(splitTag)
	for i := 0; i < 100; i++ {
		tm := time.Now().Format(splitTag)
		if tm != prevTime {
			lines, _ := u.ReadFile(logFile+"."+prevTime, 2048000)
			ok := true
			for _, timeValue := range outs {
				if !strings.Contains(lines, timeValue) {
					t.Error("log line error", prevTime, i, timeValue, u.JsonP(outs), u.JsonP(lines))
					ok = false
					break
				}
			}
			if !ok {
				break
			}
			fmt.Println("log check succeed", prevTime, len(outs))
			outs = make([]string, 0)
			prevTime = tm
		}

		fmt.Println(u.String(i) + "_" + tm)
		logger.Info(u.String(i) + "_" + tm)
		outs = append(outs, u.String(i) + "_" + tm)

		time.Sleep(100 * time.Millisecond)
		time.Sleep(3 * time.Second)
	}


	_ = os.RemoveAll(logFile + "*")
}
