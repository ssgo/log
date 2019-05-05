package log

type Logable interface {
	SetLogger(logger *Logger)
}
