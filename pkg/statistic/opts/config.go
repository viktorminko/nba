package opts

type Config struct {
	//port to check container health state
	HealthcheckPort int

	//HTTP server port
	ServerPort int

	//Path to html template file
	HTMLTemplatePath string

	//MQTT configuration
	MQTT MQTTConfig
}
