package mqtt

import (
	"context"
	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

type ClientConfig struct {
	Broker   string
	ClientID string
}

type MQTT struct {
	options *pmqtt.ClientOptions
	client  pmqtt.Client
}

func New(ctx context.Context, cfg ClientConfig) (*MQTT, error) {
	res := &MQTT{
		options: pmqtt.NewClientOptions().
			AddBroker(cfg.Broker).
			//SetCleanSession(false).
			SetClientID(cfg.ClientID),
	}

	res.client = pmqtt.NewClient(res.options)
	token := res.client.Connect()
	token.Wait()

	if token.Error() != nil {
		return nil, errors.Wrap(token.Error(), "connect to broker")
	}

	go func() {
		<-ctx.Done()
		res.client.Disconnect(1000)
	}()

	return res, nil
}
