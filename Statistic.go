package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Statistic(serverId, app, name string, startTime, endTime float64, total, failed uint, avgTime, minTime, maxTime float32, extra ...interface{}) {
	if !logger.checkLevel(INFO) {
		return
	}

	logger.log(standard.StatisticLog{
		BaseLog:   logger.getBaseLog(standard.LogTypeStatistic, extra...),
		ServerId:  serverId,
		App:       app,
		Name:      name,
		StartTime: startTime,
		EndTime:   endTime,
		Total:     total,
		Failed:    failed,
		AvgTime:   avgTime,
		MinTime:   minTime,
		MaxTime:   maxTime,
	})
}
