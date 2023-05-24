package generic

type Configuration struct {
	Version  string `env:"VERSION"`
	AppEnv   string `env:"APP_ENV"`
	LogLevel string `env:"LOG_LEVEL"`
}
