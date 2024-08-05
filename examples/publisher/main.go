package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/khan-lau/kutils/logger"
	rabbitmq "github.com/wagslane/go-rabbitmq"
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

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		glog.F("{}", err.Error())
	}
	defer publisher.Close()

	publisher.NotifyReturn(func(r rabbitmq.Return) {
		glog.D("message returned from server: {}", string(r.Body))
	})

	publisher.NotifyPublish(func(c rabbitmq.Confirmation) {
		glog.D("message confirmed from server. tag: {}, ack: {}", c.DeliveryTag, c.Ack)
	})

	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
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
				glog.E("{}", err.Error())
			}
		case <-done:
			fmt.Println("stopping publisher")
			return
		}
	}
}
