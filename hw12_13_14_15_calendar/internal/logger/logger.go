package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
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
