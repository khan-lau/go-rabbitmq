package main

import (
	"context"

	rabbitmq "github.com/khan-lau/go-rabbitmq"
	"github.com/khan-lau/kutils/logger"
)

// errorLogger is used in WithPublisherOptionsLogger to create a custom logger
// that only logs ERROR and FATAL log levels
type errorLogger struct {
	log *logger.Logger
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

func main() {
	glog := logger.LoggerInstanceOnlyConsole(int8(logger.DebugLevel))
	mylogger := &errorLogger{log: logger.LoggerInstanceOnlyConsole(int8(logger.DebugLevel))}

	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		glog.F("{}", err.Error())
	}
	defer conn.Close()

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogger(mylogger),
	)
	if err != nil {
		glog.F("{}", err.Error())
	}
	err = publisher.PublishWithContext(
		context.Background(),
		[]byte("hello, world"),
		[]string{"my_routing_key"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange("events"),
	)
	if err != nil {
		glog.F("{}", err.Error())
	}

	publisher.NotifyReturn(func(r rabbitmq.Return) {
		glog.I("message returned from server:{}", string(r.Body))
	})
}
