package rabbitmq

import (
	"fmt"
	// "log"
	klog "github.com/khan-lau/kutils/logger"
	"github.com/wagslane/go-rabbitmq/internal/logger"
)

// Logger is describes a logging structure. It can be set using
// WithPublisherOptionsLogger() or WithConsumerOptionsLogger().
type Logger logger.Logger

const loggingPrefix = "gorabbit"

type stdDebugLogger struct {
	log *klog.Logger
}

// Fatalf -
func (l stdDebugLogger) Fatalf(format string, v ...interface{}) {
	if l.log == nil {
		return
	}
	l.log.Fatal(fmt.Sprintf("%s FATAL: %s", loggingPrefix, format), v...)
}

// Errorf -
func (l stdDebugLogger) Errorf(format string, v ...interface{}) {
	if l.log == nil {
		return
	}
	l.log.Error(fmt.Sprintf("%s ERROR: %s", loggingPrefix, format), v...)
}

// Warnf -
func (l stdDebugLogger) Warnf(format string, v ...interface{}) {
	if l.log == nil {
		return
	}
	l.log.Warrn(fmt.Sprintf("%s WARN: %s", loggingPrefix, format), v...)
}

// Infof -
func (l stdDebugLogger) Infof(format string, v ...interface{}) {
	if l.log == nil {
		return
	}
	l.log.Info(fmt.Sprintf("%s INFO: %s", loggingPrefix, format), v...)
}

// Debugf -
func (l stdDebugLogger) Debugf(format string, v ...interface{}) {
	if l.log == nil {
		return
	}
	l.log.Debug(fmt.Sprintf("%s DEBUG: %s", loggingPrefix, format), v...)
}
