package main

import (
	"github.com/khan-lau/kutils/logger"
	rabbitmq "github.com/wagslane/go-rabbitmq"
)

func main() {
	glog := logger.LoggerInstanceOnlyConsole(int8(logger.DebugLevel))
	resolver := rabbitmq.NewStaticResolver(
		[]string{
			"amqp://guest:guest@host1",
			"amqp://guest:guest@host2",
			"amqp://guest:guest@host3",
		},
		false, /* shuffle */
	)

	conn, err := rabbitmq.NewClusterConn(resolver)
	if err != nil {
		glog.F("{}", err.Error())
	}
	defer conn.Close()

}
