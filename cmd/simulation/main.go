package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/viktorminko/nba/pkg/simulation/opts"
	"github.com/viktorminko/nba/pkg/simulation/service"
	"github.com/viktorminko/nba/pkg/simulation/transport/mqtt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

func readDataFile(path string) (*os.File, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}

	return file, nil
}

func initTopic(ctx context.Context, config opts.MQTTConfig) (*mqtt.TopicTransport, error) {
	client, err := mqtt.New(ctx, mqtt.ClientConfig(config.ClientConfig))
	if err != nil {
		return nil, errors.Wrap(err, "connect to mqtt")
	}

	return client.CreateTopicTransport(mqtt.TopicConfig(config.EventsTopic)), nil
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("panic: %s: %s", r, string(debug.Stack()))
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	setupGracefulShutdown(cancel)

	config := opts.ReadConfig()

	//read file with teams/players data
	file, err := readDataFile(config.DataFilePath)
	if err != nil {
		log.Fatalf("read data file: %#v", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Fatal("close file", cerr)
		}
	}()

	//create mqtt topic to publish events to
	eventsTopic, err := initTopic(ctx, config.MQTT)
	if err != nil {
		log.Fatal("init topic", err)
	}

	if err := service.Start(ctx, file, eventsTopic, config.GameDuration, config.EventDuration); err != nil {
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
