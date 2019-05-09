package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) Statistic(serverId, app, name string, startTime, endTime float64, total, failed uint, avgTime, minTime, maxTime float32, extra ...interface{}) {
	if !logger.CheckLevel(INFO) {
		return
	}

	logger.Log(logger.MakeStatisticLog(standard.LogTypeStatistic, serverId, app, name, startTime, endTime, total, failed, avgTime, minTime, maxTime, extra...))
}

func (logger *Logger) MakeStatisticLog(logType, serverId, app, name string, startTime, endTime float64, total, failed uint, avgTime, minTime, maxTime float32, extra ...interface{}) standard.StatisticLog {
	return standard.StatisticLog{
		BaseLog:   logger.MakeBaseLog(logType, extra...),
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
	}
}
