package service

import (
	"context"
	"log/slog"

	"github.com/dan-solli/homeapps/microservice/servicemesh/config"
)

type IStore interface {
	StoreService(c context.Context, s Service) error
	GetServices(c context.Context) ([]Service, error)
}

func NewBackend(cfg config.DB, f IStore, l *slog.Logger) (IStore, error) {
	return NewPgSQLRepository(cfg, l)
}
