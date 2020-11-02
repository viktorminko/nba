package opts

import (
	"time"
)

type Config struct {
	//path to file with data for teams/players
	DataFilePath string

	//duration of single game
	GameDuration time.Duration

	//how often event is fired
	EventDuration time.Duration

	//port to check if service is healthy
	HealthcheckPort int

	//configuration of MQTT queue
	MQTT MQTTConfig
}
