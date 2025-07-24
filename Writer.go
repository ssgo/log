package log

import (
	"sync"
	"time"
)

type Writer interface {
	Log([]byte)
	Run()
}

var writerRunning = false
var writerLock = sync.RWMutex{}
var writerStopChan chan bool
var writers = make([]Writer, 0)

func CheckStart() {
	writerLock.RLock()
	wr := writerRunning
	writerLock.RUnlock()
	if !wr {
		Start()
	}
}

func Start() {
	writerLock.Lock()
	defer writerLock.Unlock()
	if writerRunning {
		return
	}
	writerRunning = true
	writerStopChan = make(chan bool)
	go writerRunner()
}

func Stop() {
	writerLock.Lock()
	defer writerLock.Unlock()
	if !writerRunning {
		return
	}
	writerRunning = false
}

func Wait() {
	if writerStopChan != nil {
		<-writerStopChan
		writerStopChan = nil
	}
}

func writerRunner() {
	for {
		tmpFiles := make([]*File, len(files))
		filesLock.RLock()
		i := 0
		for _, f := range files {
			tmpFiles[i] = f
			i++
		}
		filesLock.RUnlock()
		for _, f := range tmpFiles {
			f.Run()
		}

		tmpWrites := make([]Writer, len(writers))
		writerLock.RLock()
		copy(tmpWrites, writers)
		writerLock.RUnlock()

		for _, w := range tmpWrites {
			w.Run()
		}

		writerLock.RLock()
		wr := writerRunning
		writerLock.RUnlock()

		if !wr {
			time.Sleep(5 * time.Millisecond)
			for _, f := range tmpFiles {
				f.Run()
				f.Close()
			}
			for _, w := range tmpWrites {
				w.Run()
			}
			break
		}

		// 每10毫秒Run一次
		time.Sleep(5 * time.Millisecond)
	}

	if writerStopChan != nil {
		close(writerStopChan)
	}
}
