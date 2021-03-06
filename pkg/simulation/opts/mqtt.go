package opts

import (
	"github.com/spf13/viper"
)

type MQTTClientConfig struct {
	//broker URL
	Broker string

	//ID of client
	ClientID string
}

type MQTTTopicConfig struct {
	//Name of the MQTT topic
	Name string
}

type MQTTConfig struct {
	//Topic for events
	EventsTopic MQTTTopicConfig

	//MQTT client config
	ClientConfig MQTTClientConfig
}

func initMQTTConfig(vpr *viper.Viper) MQTTConfig {

	vpr.SetDefault("MQTT_EVENTS_TOPIC_QOS", 2)
	vpr.SetDefault("MQTT_EVENTS_TOPIC", "events")

	vpr.SetDefault("MQTT_BROKER", "tcp://broker:1883")
	vpr.SetDefault("MQTT_CLIENT_ID", "nba-client")

	return MQTTConfig{
		EventsTopic: MQTTTopicConfig{
			Name: vpr.GetString("MQTT_EVENTS_TOPIC"),
		},

		ClientConfig: MQTTClientConfig{
			Broker:   vpr.GetString("MQTT_BROKER"),
			ClientID: vpr.GetString("MQTT_CLIENT_ID"),
		},
	}
}
