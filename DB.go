package log

import (
	"github.com/ssgo/standard"
	"github.com/ssgo/u"
)

func (logger *Logger) DB(dbType, dsn, query string, args []interface{}, usedTime float32, extra ...interface{}) {
	if !logger.CheckLevel(INFO) {
		return
	}
	logger.Log(logger.MakeDBLog(standard.LogTypeDb, dbType, dsn, query, args, usedTime, extra...))
}

func (logger *Logger) DBError(error, dbType, dsn, query string, args []interface{}, usedTime float32, extra ...interface{}) {
	if !logger.CheckLevel(ERROR) {
		return
	}
	logger.Log(standard.DBErrorLog{
		DBLog:    logger.MakeDBLog(standard.LogTypeDbError, dbType, dsn, query, args, usedTime, extra...),
		ErrorLog: logger.MakeErrorLog(standard.LogTypeDbError, error),
	})
}

func (logger *Logger) MakeDBLog(logType, dbType, dsn, query string, args []interface{}, usedTime float32, extra ...interface{}) standard.DBLog {
	return standard.DBLog{
		BaseLog:   logger.MakeBaseLog(logType, extra...),
		DbType:    dbType,
		Dsn:       dsn,
		Query:     query,
		QueryArgs: u.String(args),
		UsedTime:  usedTime,
	}
}
