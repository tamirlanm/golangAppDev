package logger

import "log"

type Interface interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type Logger struct{}

func New() *Logger { return &Logger{} }

func (l *Logger) Debug(message interface{}, args ...interface{}) {
	log.Println("DEBUG:", message)
}
func (l *Logger) Info(message string, args ...interface{}) {
	log.Println("INFO:", message)
}
func (l *Logger) Warn(message string, args ...interface{}) {
	log.Println("WARN:", message)
}
func (l *Logger) Error(message interface{}, args ...interface{}) {
	log.Println("ERROR:", message)
}
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	log.Fatal("FATAL:", message)
}
