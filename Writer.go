package log

import (
	"time"
)

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
		i := 0
		filesLock.RLock()
		runFiles := make([]*File, len(files))
		for _, f := range files {
			runFiles[i] = f
			i++
		}
		filesLock.RUnlock()
		for _, f := range runFiles {
			f.Run()
		}
		for _, w := range writers {
			w.Run()
		}

		if !writerRunning {
			time.Sleep(10 * time.Millisecond)
			for _, f := range runFiles {
				f.Run()
				f.Close()
			}
			for _, w := range writers {
				w.Run()
			}
			break
		}

		// 每10毫秒Run一次
		time.Sleep(10 * time.Millisecond)
	}
	writerStopChan <- true
}
