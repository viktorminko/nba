package opts

import (
	"github.com/spf13/viper"
	"time"
)

func ReadConfig() *Config {
	vpr := viper.New()
	vpr.AutomaticEnv()

	vpr.SetDefault("HEALTHCHECK_PORT", 8888)
	vpr.SetDefault("DATA_FILE_PATH", "/players.json")
	vpr.SetDefault("GAME_DURATION", 240*time.Second)
	vpr.SetDefault("EVENT_DURATION", 5*time.Second)

	return &Config{
		HealthcheckPort: vpr.GetInt("HEALTHCHECK_PORT"),
		DataFilePath:    vpr.GetString("DATA_FILE_PATH"),
		GameDuration:    vpr.GetDuration("GAME_DURATION"),
		EventDuration:   vpr.GetDuration("EVENT_DURATION"),
		MQTT:            initMQTTConfig(vpr),
	}
}
