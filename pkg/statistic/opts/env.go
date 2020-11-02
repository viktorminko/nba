package opts

import (
	"github.com/spf13/viper"
)

func ReadConfig() *Config {
	vpr := viper.New()
	vpr.AutomaticEnv()

	vpr.SetDefault("HEALTHCHECK_PORT", 8888)
	vpr.SetDefault("SERVICE_PORT", 8080)
	vpr.SetDefault("HTML_TEMPLATE_FILE_PATH", "/layout.html")

	return &Config{
		HealthcheckPort:  vpr.GetInt("HEALTHCHECK_PORT"),
		ServerPort:       vpr.GetInt("SERVICE_PORT"),
		HTMLTemplatePath: vpr.GetString("HTML_TEMPLATE_FILE_PATH"),
		MQTT:             initMQTTConfig(vpr),
	}
}
