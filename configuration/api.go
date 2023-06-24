package configuration

type Api struct {
	Port         string   `env:"PORT" envDefault:"8080"`
	InternalPort string   `env:"INTERNAL_PORT" envDefault:"2112"`
	Postgres     Postgres `envPrefix:"POSTGRES_"`
}

type Postgres struct {
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Database string `env:"DATABASE"`
	Sslmode  string `env:"SSL_MODE"`
}
