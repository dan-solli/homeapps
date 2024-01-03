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

func NewServiceMeshConfig() ServiceMeshConfig {
	viper.SetEnvPrefix("MS_SM")
	viper.AutomaticEnv()

	return ServiceMeshConfig{
		tls:      viper.GetBool("TLS"),
		certFile: viper.GetString("CERTFILE"),
		keyFile:  viper.GetString("KEYFILE"),
		grpcport: viper.GetInt("GRPC_PORT"),
		httpport: viper.GetInt("HTTP_PORT"),
	}
}

func (c *ServiceMeshConfig) WithTLS(tls bool) ServiceMeshConfig {
	c.tls = tls
	return *c
}

func (c *ServiceMeshConfig) WithCertFile(certFile string) ServiceMeshConfig {
	c.certFile = certFile
	return *c
}

func (c *ServiceMeshConfig) WithKeyFile(keyFile string) ServiceMeshConfig {
	c.keyFile = keyFile
	return *c
}

func (c *ServiceMeshConfig) WithGrpcPort(grpcport int) ServiceMeshConfig {
	c.grpcport = grpcport
	return *c
}

func (c *ServiceMeshConfig) WithHttpPort(httpport int) ServiceMeshConfig {
	c.httpport = httpport
	return *c
}
