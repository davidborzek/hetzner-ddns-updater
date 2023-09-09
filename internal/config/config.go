package config

import (
	"os"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var (
	envOpts = env.Options{
		Prefix: "HDU_",
	}
)

type Config struct {
	Address        string        `env:"ADDRESS" default:":8080"`
	ApiToken       string        `env:"API_TOKEN" validate:"required"`
	RecordID       string        `env:"RECORD_ID" validate:"required"`
	ZoneID         string        `env:"ZONE_ID" validate:"required"`
	RecordName     string        `env:"RECORD_NAME" default:"@"`
	RecordTTL      uint          `env:"RECORD_TTL" default:"60"`
	Interval       time.Duration `env:"INTERVAL" default:"5m"`
	MetricsEnabled bool          `env:"METRICS_ENABLED" default:"false"`
	MetricsToken   string        `env:"METRICS_TOKEN"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := defaults.Set(&cfg); err != nil {
		return nil, err
	}

	if err := godotenv.Load("hetzner-ddns-updater.config"); err != nil && err == os.ErrNotExist {
		return nil, err
	}

	if err := env.ParseWithOptions(&cfg, envOpts); err != nil {
		return nil, err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
