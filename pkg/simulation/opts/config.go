package opts

import (
	"time"
)

type Config struct {
	DataFilePath  string
	GameDuration  time.Duration
	EventDuration time.Duration

	HealthcheckPort int

	MQTT MQTTConfig
}
