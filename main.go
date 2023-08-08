package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/creasty/defaults"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPathEnvVarName = "APP_CONFIG_PATH"
)

const (
	MySQLDriver = "mysql"
)

type config struct {
	Http     httpConfig     `yaml:"http"`
	Database databaseConfig `yaml:"database"`
	Log      logConfig      `yaml:"log"`
}

type httpConfig struct {
	GinMode string `yaml:"gin_mode" env:"GIN_MODE" env-default:"debug"`
	Port    int    `yaml:"port" env:"PORT" env-default:"80"`
}

func (c httpConfig) Validate() error {
	if c.GinMode != gin.DebugMode && c.GinMode != gin.ReleaseMode && c.GinMode != gin.TestMode {
		return errors.New("invalid `gin_mode`")
	}
	if c.Port < 1 || c.Port > 65536 {
		return errors.New("invalid `port`")
	}
	return nil
}

type databaseConfig struct {
	Driver          string        `yaml:"driver"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	Address         string        `yaml:"address"`
	DatabaseName    string        `yaml:"database_name"`
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"10"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"20"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env-default:"30s"`
}

func (c databaseConfig) Validate() error {
	if c.Driver != MySQLDriver {
		return errors.New("invalid `driver`")
	}
	if c.Username == "" {
		return errors.New("missing `username`")
	}
	if c.Address == "" {
		return errors.New("missing `address`")
	}
	if c.DatabaseName == "" {
		return errors.New("missing `database_name`")
	}

	return nil
}

type logConfig struct {
	DevLogMode bool   `yaml:"dev_log_mode" default:"true"`
	LogLevel   string `yaml:"log_level" default:"info"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

func (c *logConfig) Validate() error {
	if c.Filename == "" {
		return errors.New("missing `filename`")
	}
	return nil
}

// Validate application configuration
func (c config) Validate() error {
	if err := c.Http.Validate(); err != nil {
		return fmt.Errorf("http.(%s)", err.Error())
	}

	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("log.(%s)", err.Error())
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database.(%s)", err.Error())
	}

	return nil
}

// SetDefaultValues set default values for the configuration
func (c *config) SetDefaultValues() error {
	return defaults.Set(c)
}

func loadConfig(configFile string) (*config, error) {
	var cfg config

	err := cleanenv.ReadConfig(configFile, &cfg)
	if err != nil {
		return nil, err
	}

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	return &cfg, nil
}

func newConfig() (*config, error) {
	if os.Getenv(configPathEnvVarName) == "" {
		return nil, fmt.Errorf("missing `%s` environment variable", configPathEnvVarName)
	}

	config, err := loadConfig(os.Getenv(configPathEnvVarName))
	if err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	fmt.Println("try git ci-cd")
}
