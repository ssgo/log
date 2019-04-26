package log

import (
	"github.com/ssgo/standard"
)

func (logger *Logger) DB(dbType, dsn, query string, args []interface{}, usedTime float32, extra ...interface{}) {
	if !logger.checkLevel(INFO) {
		return
	}
	logger.log(logger.getDBLog(standard.LogTypeDb, dbType, dsn, query, args, usedTime, extra...))
}

func (logger *Logger) DBError(error, dbType, dsn, query string, args []interface{}, usedTime float32, extra ...interface{}) {
	if !logger.checkLevel(ERROR) {
		return
	}
	logger.log(standard.DBErrorLog{
		DBLog:    logger.getDBLog(standard.LogTypeDbError, dbType, dsn, query, args, usedTime, extra...),
		ErrorLog: standard.ErrorLog{Error: error},
	})
}

func (logger *Logger) getDBLog(logType, dbType, dsn, query string, args []interface{}, usedTime float32, extra ...interface{}) standard.DBLog {
	return standard.DBLog{
		BaseLog:  logger.getBaseLog(logType, extra...),
		DbType:   dbType,
		Dsn:      dsn,
		Query:    query,
		Args:     args,
		UsedTime: usedTime,
	}
}
