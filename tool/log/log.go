package log

import (
	"fmt"
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

type Option func(*Logger)

func WithLevel(level logrus.Level) Option {
	return func(l *Logger) {
		l.SetLevel(level)
	}
}

type MultiWriter struct {
	writers []io.Writer
}

func NewMultiWriter(writers ...io.Writer) *MultiWriter {
	return &MultiWriter{
		writers: writers,
	}
}

func (w *MultiWriter) Write(p []byte) (n int, err error) {
	for i := 0; i < len(w.writers)-1; i++ {
		w.writers[i].Write(p)
	}
	if len(w.writers) <= 0 {
		return 0, nil
	}
	return w.writers[len(w.writers)-1].Write(p)
}
func WithConsole() Option {
	return func(l *Logger) {
		l.SetOutput(NewMultiWriter(
			l.Out,
			os.Stdout,
		))
	}
}

func WithFilePath(path string) Option {
	return func(l *Logger) {
		changePath := func(writer io.Writer, path string) bool {
			if out, ok := writer.(*lumberjack.Logger); ok {
				out.Filename = fmt.Sprintf("%s/%s", path, out.Filename)
				return true
			}
			return false
		}
		ok := changePath(l.Out, path)
		if !ok {
			if out, ok := l.Out.(*MultiWriter); ok {
				for _, w := range out.writers {
					changePath(w, path)
				}
			}
		}
	}
}

func WithFormater(formatter logrus.Formatter) Option {
	return func(l *Logger) {
		l.SetFormatter(formatter)
	}
}

func NewAppLogger(app string, options ...Option) *Logger {
	logger := &Logger{
		Logger: NewLogger(app),
	}
	for _, o := range options {
		o(logger)
	}
	return logger
}
