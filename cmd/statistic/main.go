package main

import (
	"context"
	"github.com/pkg/errors"
	htmlfront "github.com/viktorminko/nba/pkg/statistic/frontend/html"
	"github.com/viktorminko/nba/pkg/statistic/opts"
	"github.com/viktorminko/nba/pkg/statistic/service"
	"github.com/viktorminko/nba/pkg/statistic/subscriber/mqtt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func initTopics(ctx context.Context, config opts.MQTTConfig) (*mqtt.TopicSubscriber, error) {
	client, err := mqtt.New(ctx, mqtt.ClientConfig(config.ClientConfig))
	if err != nil {
		return nil, errors.Wrap(err, "connect to mqtt")
	}

	subEvents, err := client.CreateTopicSubscriber(ctx, mqtt.TopicConfig(config.EventsTopic))
	if err != nil {
		return nil, errors.Wrap(err, "create topic subsciber")
	}

	return subEvents, nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("panic: %s: %s", r, string(debug.Stack()))
		}
	}()
	// make context
	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	config := opts.ReadConfig()

	topicEvents, err := initTopics(ctx, config.MQTT)
	if err != nil {
		log.Fatal("init topic", err)
	}

	fe, err := htmlfront.New(config.HTMLTemplatePath)
	if err != nil {
		log.Fatal("build html frontend", err)
	}

	if err := service.Start(ctx, topicEvents, config.ServerPort, fe); err != nil {
		log.Fatal("error starting service", err)
	}
}

func setupGracefulShutdown(stop func()) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Println("Got Interrupt signal")
		stop()
	}()
}
