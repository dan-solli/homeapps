package main

import (
	"github.com/spf13/viper"
)

type ServiceMeshConfig struct {
	tls      bool
	certFile string
	keyFile  string
	grpcport int
	httpport int
}

// TODO: Need to fake the config. But that might be easy by not calling this func. Except if it contains more stuff.
func NewServiceMeshConfig() (ServiceMeshConfig, error) {
	viper.SetEnvPrefix("MS_SM")
	viper.AutomaticEnv()

	return ServiceMeshConfig{
		tls:      viper.GetBool("TLS"),
		certFile: viper.GetString("CERTFILE"),
		keyFile:  viper.GetString("KEYFILE"),
		grpcport: viper.GetInt("GRPC_PORT"),
		httpport: viper.GetInt("HTTP_PORT"),
	}, nil
}
