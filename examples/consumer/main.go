package main

import (
	"fmt"

	"os"
	"os/signal"
	"syscall"

	rabbitmq "github.com/khan-lau/go-rabbitmq"
	"github.com/khan-lau/kutils/logger"
)

func main() {
	glog := logger.LoggerInstanceOnlyConsole(int8(logger.DebugLevel))
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		glog.F("{}", err.Error())
	}
	defer conn.Close()

	consumer, err := rabbitmq.NewConsumer(
		conn,
		"my_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		glog.F("{}", err.Error())
	}

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("awaiting signal")
		sig := <-sigs

		fmt.Println()
		fmt.Println(sig)
		fmt.Println("stopping consumer")

		consumer.Close()
	}()

	// block main thread - wait for shutdown signal
	err = consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		glog.I("consumed: {}", string(d.Body))

		// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
		return rabbitmq.Ack
	})
	if err != nil {
		glog.F("{}", err.Error())
	}
}
