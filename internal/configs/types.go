package configs

type Config struct {
	Server        Server
	Database      Database
	FositeConfigs FositeConfigs
}

type Server struct {
	Port string `env:"PORT" envDefault:"8080"`
}

type Database struct {
	Username string `env:"DB_USERNAME,required"`
	Password string `env:"DB_PASSWORD_FILE,file,required"`
	Name     string `env:"DB_NAME,required"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

type FositeConfigs struct {
	Secret string `env:"FOSITE_SECRET_FILE,file,required"`
}
