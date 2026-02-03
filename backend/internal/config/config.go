package config

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

type Config struct {
	BaseConfig     `mapstructure:",squash"`
	PostgresConfig `mapstructure:",squash"`
	RedisConfig    `mapstructure:",squash"`
	HTTPConfig     `mapstructure:",squash"`
	ParserConfig   `mapstructure:",squash"`
}

type BaseConfig struct {
	Environment string
	LogMode     string `mapstructure:"LOG_MODE"`
}

type PostgresConfig struct {
	PgHost     string `mapstructure:"PG_HOST"`
	PgPort     string `mapstructure:"PG_PORT"`
	PgUser     string `mapstructure:"PG_USER"`
	PgPassword string `mapstructure:"PG_PASSWORD"`
	PgDBName   string `mapstructure:"PG_DB_NAME"`
}

func (c *PostgresConfig) GetHost() string     { return c.PgHost }
func (c *PostgresConfig) GetPort() string     { return c.PgPort }
func (c *PostgresConfig) GetUser() string     { return c.PgUser }
func (c *PostgresConfig) GetPassword() string { return c.PgPassword }
func (c *PostgresConfig) GetDBName() string   { return c.PgDBName }

type RedisConfig struct {
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

type HTTPConfig struct {
	HTTPPort       string `mapstructure:"HTTP_PORT"`
	GinMode        string `mapstructure:"GIN_MODE"`
	RateLimitRPS   int    `mapstructure:"RATE_LIMIT_RPS"`
	RateLimitBurst int    `mapstructure:"RATE_LIMIT_BURST"`
}

type ParserConfig struct {
	ScrapeInterval time.Duration `mapstructure:"SCRAPE_INTERVAL"`
	ScrapeWorkers  int           `mapstructure:"SCRAPE_WORKERS"`
}

func LoadFromFlags(cfg *Config) error {
	var envFile string
	flag.StringVar(&envFile, "env", "", "path to .env file")
	flag.Parse()

	return Load(cfg, envFile)
}

func Load(cfg *Config, envFile string) error {
	v := viper.New()

	setDefaults(v)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if envFile != "" {
		v.SetConfigFile(envFile)
	} else {
		v.SetConfigName(".env")
		v.SetConfigType("env")
		v.AddConfigPath(".")
		v.AddConfigPath("../..")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("reading config: %w", err)
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshaling config: %w", err)
	}

	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("LOG_MODE", "dev")

	v.SetDefault("PG_HOST", "localhost")
	v.SetDefault("PG_PORT", "5432")
	v.SetDefault("PG_USER", "postgres")
	v.SetDefault("PG_PASSWORD", "postgres")
	v.SetDefault("PG_DB_NAME", "store_scraper")

	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", "6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	v.SetDefault("HTTP_PORT", "8080")
	v.SetDefault("GIN_MODE", "debug")
	v.SetDefault("RATE_LIMIT_RPS", 100)
	v.SetDefault("RATE_LIMIT_BURST", 200)

	v.SetDefault("SCRAPE_INTERVAL", 10*time.Minute)
	v.SetDefault("SCRAPE_WORKERS", 5)
}

func (c *BaseConfig) IsDevEnv() bool {
	return c.Environment == EnvDev
}

func (c *BaseConfig) IsProdEnv() bool {
	return c.Environment == EnvProd
}
