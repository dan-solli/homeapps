package config

import (
	"log/slog"

	"github.com/iamolegga/enviper"

	"github.com/spf13/viper"
)

/*
type Config struct {
	Rest_port           int    `mapstructure:"rest_port"`
	Metrics_path        string `mapstructure:"metrics_path"`
	Health_path         string `mapstructure:"health_path"`
	Grpc_port           int    `mapstructure:"grpc_port"`
	Tls                 bool   `mapstructure:"tls"`
	Keyfile             string `mapstructure:"keyfile"`
	Certfile            string `mapstructure:"certfile"`
	Eventbroker_address string `mapstructure:"eventbroker_address"`
	Db_host             string `mapstructure:"DB_HOST"`
	Db_port             int    `mapstructure:"db_port"`
	Db_user             string `mapstructure:"db_user"`
	Db_pass             string `mapstructure:"db_pass"`
}
*/

type Http struct {
	Rest_port    int    `mapstructure:"REST_PORT"`
	Metrics_path string `mapstructure:"metrics_path"`
	Health_path  string `mapstructure:"HEALTH_PATH"`
}

type GRpc struct {
	Grpc_port int    `mapstructure:"grpc_port"`
	Tls       bool   `mapstructure:"tls"`
	Keyfile   string `mapstructure:"keyfile"`
	Certfile  string `mapstructure:"certfile"`
}

type DB struct {
	Db_host string `mapstructure:"DB_HOST"`
	Db_port int    `mapstructure:"db_port"`
	Db_user string `mapstructure:"db_user"`
	Db_pass string `mapstructure:"db_pass"`
	Db_name string `mapstructure:"db_name"`
}

type Server struct {
	Http Http
	GRpc GRpc
	DB   DB
}

type Client struct {
	Eventbroker_address string `mapstructure:"eventbroker_address"`
}

type Config struct {
	Server Server
	Client Client
}

func (c *Http) Port() int {
	return c.Rest_port
}

var e = enviper.New(viper.New())

func NewConfig(log *slog.Logger) *Config {
	e.SetConfigFile(".env")
	e.SetConfigType("env")
	e.AddConfigPath(".")
	e.AutomaticEnv()

	cfg := Config{
		Server: Server{
			Http: Http{},
			GRpc: GRpc{},
			DB:   DB{},
		},
		Client: Client{
			Eventbroker_address: "",
		},
	}

	log.Info("Empty config", "cfg", cfg)

	if err := e.ReadInConfig(); err != nil {
		log.Error("Failed to read config file", "err", err)
		panic(err)
	}

	return &cfg
}

func Viper() *viper.Viper {
	return e.Viper
}
