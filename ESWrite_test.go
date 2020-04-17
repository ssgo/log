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
	logger.Info("Test", "level", 2)
	logger.Warning("Test", "level", 3)
	logger.Error("Test", "level", 4)
	log.Stop()
	log.Wait()
}
