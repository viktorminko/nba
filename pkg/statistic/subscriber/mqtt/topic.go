package mqtt

import (
	"context"
	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
)

type TopicConfig struct {
	Name string
}

//see Transporter interface
type TopicSubscriber struct {
	dataCh      <-chan []byte
	healthCheck func() error
}

func (s *TopicSubscriber) Subscribe() <-chan []byte {
	return s.dataCh
}

func (s *TopicSubscriber) HealthCheck() error {
	return s.healthCheck()
}

func (m *MQTT) CreateTopicSubscriber(ctx context.Context, cfg TopicConfig) (*TopicSubscriber, error) {

	dataCh := make(chan []byte)

	token := m.client.Subscribe(cfg.Name, 2, func(c pmqtt.Client, mqttMsg pmqtt.Message) {
		dataCh <- mqttMsg.Payload()
	})
	token.Wait()

	if token.Error() != nil {
		return nil, errors.Wrap(token.Error(), "subscribe to topic")
	}

	return &TopicSubscriber{
		dataCh: dataCh,
		healthCheck: func() error {
			//no error if connected or trying to reconnect
			if !m.client.IsConnected() {
				return errors.New("client status is not connected or reconnecting")
			}

			return nil
		},
	}, nil
}
