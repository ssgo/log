package log

import "time"

type Writer interface {
	Log([]byte)
	Run()
}

var writerRunning = false
var writerStopChan = make(chan bool)
var writers = make([]Writer, 0)

func Start() {
	writerRunning = true
	go writerRunner()
}

func Stop() {
	writerRunning = false
}

func Wait() {
	<-writerStopChan
}

func writerRunner() {
	for {
		for _, w := range writers {
			w.Run()
		}

		if !writerRunning {
			time.Sleep(100 * time.Millisecond)
			for _, w := range writers {
				w.Run()
			}
			break
		}

		// 每100毫秒Run一次
		time.Sleep(100 * time.Millisecond)
	}
	writerStopChan <- true
}
