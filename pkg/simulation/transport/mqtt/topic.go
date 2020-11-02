package mqtt

import (
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
)

type TopicConfig struct {
	Name string
}

type TopicTransport struct {
	sender      func(r io.Reader) error
	healthCheck func() error
}

func (s *TopicTransport) Transport(r io.Reader) error {
	return s.sender(r)
}

func (s *TopicTransport) HealthCheck() error {
	return s.healthCheck()
}

func (m *MQTT) CreateTopicTransport(cfg TopicConfig) *TopicTransport {
	return &TopicTransport{
		func(r io.Reader) (err error) {
			b, err := ioutil.ReadAll(r)
			if err != nil {
				return errors.Wrap(err, "failed to read serialized data")
			}

			token := m.client.Publish(cfg.Name, 2, false, b)
			token.Wait()

			if token.Error() != nil {
				return errors.Wrapf(
					token.Error(),
					"publish message to topic: %v, QoS: 2", cfg.Name,
				)
			}

			return nil
		},
		func() error {
			//no error if connected or trying to reconnect
			if !m.client.IsConnected() {
				return errors.New("client status is not connected or reconnecting")
			}

			return nil
		},
	}
}
