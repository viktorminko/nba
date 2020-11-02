package opts

type Config struct {
	HealthcheckPort int
	ServerPort      int

	HTMLTemplatePath string

	MQTT MQTTConfig
}
