package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logg *logrus.Logger
}

var errLogCreateStr = "fail in creating logger: %w"

func New(logLevel, path string) (*Logger, error) {
	l := logrus.New()

	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf(errLogCreateStr, err)
	}
	l.SetLevel(level)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644) //nolint
	if err != nil {
		return nil, fmt.Errorf(errLogCreateStr, err)
	}
	l.SetOutput(file)

	return &Logger{l}, nil
}

func (l *Logger) Info(args ...interface{}) {
	l.logg.Info(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logg.Error(args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.logg.Debug(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logg.Warn(args...)
}
