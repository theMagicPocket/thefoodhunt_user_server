package log

import (
	"log"
	"os"
)

type AppLogger interface {
	Info(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
}

type applogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func New() AppLogger {
	return &applogger{
		infoLogger:  log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *applogger) Info(args ...interface{}) {
	l.infoLogger.Println(args...)
}

func (l *applogger) Fatal(args ...interface{}) {
	l.errorLogger.Fatal(args...)
}

func (l *applogger) Error(args ...interface{}) {
	l.errorLogger.Println(args...)
}
