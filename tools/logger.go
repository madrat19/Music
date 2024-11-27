package tools

import "log"

var Logger logger

type logger struct {
	level    string
	infoLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger
}

func (l *logger) Info(message string) {
	if l.level == "info" {
		l.infoLog.Print(message)
	}
}

func (l *logger) Error(message string, err error) {
	if l.level == "info" || l.level == "error" {
		l.errorLog.Print(message, err)
	}
}

func (l *logger) Fatal(message string, err error) {
	l.fatalLog.Fatal(message, err)
}

func InitLogger(level string, infoLog *log.Logger, errorLog *log.Logger, fatalLog *log.Logger) {

	Logger = logger{
		level:    level,
		infoLog:  infoLog,
		errorLog: errorLog,
		fatalLog: fatalLog,
	}
}
