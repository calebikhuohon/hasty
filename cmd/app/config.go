package main

type Config struct {
	MySQL struct {
		User   string `envconfig:"MYSQL_USER" required:"true"`
		Passwd string `envconfig:"MYSQL_PASSWORD" required:"true"`
		Host   string `envconfig:"MYSQL_HOST" required:"true"`
		DBName string `envconfig:"MYSQL_DATABASE" required:"true"`
	}
	Ports struct {
		HTTP    string `envconfig:"HTTP_PORT" default:":8080"`
	}
	Redis struct {
		Address  string `envconfig:"REDIS_ADDRESS" required:"true"`
		Password string `envconfig:"REDIS_PASSWORD" required:"true"`
	}
	AppEnv        string `envconfig:"APP_ENV" default:"dev"`
	JobMaxRetries int    `envconfig:"JOB_MAX_RETRIES" default:"2"`
}

func (c Config) Validate() error {
	return nil
}
