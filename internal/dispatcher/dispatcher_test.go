package dispatcher

import (
	"testing"
	"time"

	klogger "github.com/khan-lau/kutils/logger"
)

var (
	mylogger = &errorLogger{log: klogger.LoggerInstanceOnlyConsole(int8(klogger.DebugLevel))}
)

func init() {
	/* load test data */
}

type errorLogger struct {
	log *klogger.Logger
}

func (l errorLogger) Fatalf(format string, v ...interface{}) {
	l.log.F("mylogger: "+format, v...)
}

func (l errorLogger) Errorf(format string, v ...interface{}) {
	l.log.E("mylogger: "+format, v...)
}

func (l errorLogger) Warnf(format string, v ...interface{}) {
}

func (l errorLogger) Infof(format string, v ...interface{}) {
}

func (l errorLogger) Debugf(format string, v ...interface{}) {
}

func TestNewDispatcher(t *testing.T) {

	d := NewDispatcher(mylogger)
	if d.subscribers == nil {
		t.Error("Dispatcher subscribers is nil")
	}
	if d.subscribersMu == nil {
		t.Error("Dispatcher subscribersMu is nil")
	}
}

func TestAddSubscriber(t *testing.T) {
	d := NewDispatcher(mylogger)
	d.AddSubscriber()
	if len(d.subscribers) != 1 {
		t.Error("Dispatcher subscribers length is not 1")
	}
}

func TestCloseSubscriber(t *testing.T) {
	d := NewDispatcher(mylogger)
	_, closeCh := d.AddSubscriber()
	close(closeCh)
	time.Sleep(time.Millisecond)
	if len(d.subscribers) != 0 {
		t.Error("Dispatcher subscribers length is not 0")
	}
}
