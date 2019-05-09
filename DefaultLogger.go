package log

import "github.com/ssgo/config"

var DefaultLogger *Logger

func init() {
	conf := Config{}
	config.LoadConfig("Log", &conf)
	if conf.Level == "" {
		conf.Level = "info"
	}
	if conf.Truncations == nil {
		conf.Truncations = []string{"github.com/", "golang.org/", "/ssgo/"}
	}
	if conf.Sensitive == nil {
		conf.Sensitive = []string{"password", "secure", "token", "accessToken"}
	}
	if conf.RegexSensitive == nil {
		conf.RegexSensitive = []string{}
	}
	if conf.SensitiveRule == nil {
		conf.SensitiveRule = []string{"12:4*4", "11:3*4", "7:2*2", "3:1*1", "2:1*0"}
	}

	DefaultLogger = NewLogger(conf)
}

func New(traceId string) *Logger {
	newLogger := *DefaultLogger
	newLogger.traceId = traceId
	return &newLogger
}
