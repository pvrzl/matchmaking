package config

import (
	"app/pkg/database"
	"app/pkg/encryption"
	"app/pkg/monitoring"

	"github.com/joho/godotenv"
)

const (
	DEVELOPMENT = "development"
	PRODUCTION  = "production"
)

type AppConfig struct {
	BuildTag  string
	Version   string
	BuildDate string
	Commit    string
	Branch    string
}

type Config struct {
	AppConfig AppConfig
	Name      string
	Env       string
	PORT      string
	Database  struct {
		Read  database.Config
		Write database.Config
	}
	Redis struct {
		Host     string
		Port     string
		Password string
	}
	APM       monitoring.Config
	EncConfig encryption.EncConfig
}

var (
	cfg Config
)

func (c *Config) readConfig() {
	godotenv.Load()
	c.Name = GetEnvString("NAME", "app")
	c.Env = GetEnvString("ENV", DEVELOPMENT)
	c.PORT = GetEnvString("APP_PORT", ":8080")

	// db read
	c.Database.Read.Host = GetEnvString("DB_READ_HOST", "db")
	c.Database.Read.UserName = GetEnvString("DB_READ_USER", "postgres")
	c.Database.Read.Password = GetEnvString("DB_READ_PASSWORD", "root")
	c.Database.Read.Name = GetEnvString("DB_READ_NAME", "db_dev")
	c.Database.Read.Port = GetEnvString("DB_READ_PORT", "5432")
	c.Database.Read.AppName = c.Name
	c.Database.Read.Extras = GetEnvString("DB_READ_EXTRAS", "")

	// db write
	c.Database.Write.Host = GetEnvString("DB_WRITE_HOST", "db")
	c.Database.Write.UserName = GetEnvString("DB_WRITE_USER", "postgres")
	c.Database.Write.Password = GetEnvString("DB_WRITE_PASSWORD", "root")
	c.Database.Write.Name = GetEnvString("DB_WRITE_NAME", "db_dev")
	c.Database.Write.Port = GetEnvString("DB_WRITE_PORT", "5432")
	c.Database.Write.AppName = c.Name
	c.Database.Write.Extras = GetEnvString("DB_WRITE_EXTRAS", "")

	// redis
	c.Redis.Host = GetEnvString("REDIS_HOST", "localhost")
	c.Redis.Port = GetEnvString("REDIS_PORT", "6379")
	c.Redis.Password = GetEnvString("REDIS_PASSWORD", "")

	// APM
	c.APM.Name = GetEnvString("APM_NAME", "local-dev")
	c.APM.License = GetEnvString("APM_LICENSE", "license")
	c.APM.Enabled = GetEnvBool("APM_ENABLED", false)

	// encryption
	c.EncConfig.Key = GetEnvString("AES_KEY", "")
	c.EncConfig.XORKey = int64(GetEnvInt("XOR_KEY", 1234356789))
}

func init() {
	cfg.readConfig()
}

func Get() *Config {
	return &cfg
}

func SetAppConfig(acfg AppConfig) {
	cfg.AppConfig = acfg
}
