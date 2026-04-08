package config

type Config struct {
	Port   int `env:"PORT"`
	AppEnv int `env:"APP_ENV"`

	DatabaseURL string `env:"DB"`

	JWTSecret string `env:"JWT_SECRET"`

	AzureStorageURI string `env:"AZURE_STORAGE_URI"`
	AzureContainer  string `env:"AZURE_CONTAINER"`

	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
}

func newConfig() *Config {
	return &Config{
		Port: 
	} 
}
