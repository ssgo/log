package log_test

import (
	"github.com/ssgo/log"
	"testing"
)


func TestES(t *testing.T) {
	log.Start()

	logger := log.NewLogger(log.Config{})
	//logger := log.NewLogger(log.Config{File: "es://elastic:VqyJOSrh8iP1JvZHEVqZ@localhost:9200/bb?timeout=10s"})
	logger.Debug("Test", "level", 1)
	logger.Info("Test", "level")
	logger.Warning("Test", "level")
	logger.Error("Test", "level", 4)
	log.Stop()
	log.Wait()

	//imageName := "dev.xue.fun:5000/ssgo/deploy:0.1.1#33"
	//a := strings.Split(imageName, "/")
	//imageName = a[len(a)-1]
	//imageName = strings.SplitN(imageName, ":", 2)[0]
	//imageName = strings.SplitN(imageName, "#", 2)[0]
	//fmt.Println(" ^^^", imageName)
}
