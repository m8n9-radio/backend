package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger(logLevel string) *Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		level = logrus.InfoLevel
		logger.WithError(err).Warn("Invalid log level provided, defaulting to info")
	}
	logger.SetLevel(level)

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) GetLevel() logrus.Level {
	return l.logger.Level
}

func (l *Logger) GetLogger() *logrus.Logger {
	return l.logger
}

func (l *Logger) Trace(args ...interface{}) {
	l.logger.Trace(args...)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

func (l *Logger) WithError(err error) *logrus.Entry {
	return l.logger.WithError(err)
}

// WithComponent returns a logger entry with component field.
func (l *Logger) WithComponent(component string) *logrus.Entry {
	return l.logger.WithField("component", component)
}

// WithOperation returns a logger entry with operation field.
func (l *Logger) WithOperation(operation string) *logrus.Entry {
	return l.logger.WithField("operation", operation)
}

// WithContext returns a logger entry with component and operation fields.
func (l *Logger) WithContext(component, operation string) *logrus.Entry {
	return l.logger.WithFields(logrus.Fields{
		"component": component,
		"operation": operation,
	})
}

// WithCorrelationID returns a logger entry with correlation_id field.
func (l *Logger) WithCorrelationID(correlationID string) *logrus.Entry {
	return l.logger.WithField("correlation_id", correlationID)
}
