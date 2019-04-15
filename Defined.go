package log

import "time"

type BaseLog struct {
	LogType  string
	LogTime  time.Time
	LogLevel string
	Traces   string
	Extra    map[string]interface{}
}

type RequestLog struct {
	BaseLog
	App                string
	Node               string
	ClientIp           string
	FromApp            string
	FromNode           string
	ClientId           string
	SessionId          string
	RequestId          string
	Host               string
	Scheme             string
	Proto              string
	AuthLevel          int
	Priority           int
	Method             string
	Path               string
	RequestHeaders     map[string]string
	RequestData        map[string]interface{}
	UsedTime           float32
	ResponseCode       int
	ResponseHeaders    map[string]string
	ResponseDataLength uint
	ResponseData       interface{}
}
