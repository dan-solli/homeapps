package service

import (
	"context"
)

type IStore interface {
	StoreService(c context.Context, s Service) error
}

type DBConfig struct {
	host string
	port int
	user string
	pass string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		host: "localhost",
		port: 5432,
		user: "servicemesh",
		pass: "pwd_servicemesh",
	}
}

func NewBackend(cfg DBConfig, f IStore) (IStore, error) {
	return NewPgSQLRepository(cfg)
}

func (c *DBConfig) WithHost(host string) DBConfig {
	c.host = host
	return *c
}

func (c *DBConfig) WithPort(port int) DBConfig {
	c.port = port
	return *c
}

func (c *DBConfig) WithUser(user string) DBConfig {
	c.user = user
	return *c
}

func (c *DBConfig) WithPass(pass string) DBConfig {
	c.pass = pass
	return *c
}
