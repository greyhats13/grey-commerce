// Path: grey-user/pkg/logger/logger.go

package logger

type Field struct {
	Key   string
	Value interface{}
}

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, err error, fields ...Field)
}
