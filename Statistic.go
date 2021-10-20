package log

import (
	"github.com/ssgo/standard"
	"time"
)

func (logger *Logger) Statistic(serverId, app, name string, startTime, endTime time.Time, times, failed uint, avg, min, max float64, extra ...interface{}) {
	if !logger.CheckLevel(INFO) {
		return
	}

	logger.Log(logger.MakeStatisticLog(standard.LogTypeStatistic, serverId, app, name, startTime, endTime, times, failed, avg, min, max, extra...))
}

func (logger *Logger) MakeStatisticLog(logType, serverId, app, name string, startTime, endTime time.Time, times, failed uint, avg, min, max float64, extra ...interface{}) standard.StatisticLog {
	return standard.StatisticLog{
		BaseLog:   logger.MakeBaseLog(logType, extra...),
		ServerId:  serverId,
		App:       app,
		Name:      name,
		StartTime: MakeLogTime(startTime),
		EndTime:   MakeLogTime(endTime),
		Times:     times,
		Failed:    failed,
		Avg:       avg,
		Min:       min,
		Max:       max,
	}
}
