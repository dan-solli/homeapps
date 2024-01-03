package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/spf13/viper"
)

type PgSQLRepository struct {
	lock *sync.RWMutex
	db   *sql.DB
	l    *slog.Logger
}

func NewPgSQLRepository(d DBConfig) (*PgSQLRepository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetInt("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASS"),
		viper.GetString("DB_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return &PgSQLRepository{
		lock: &sync.RWMutex{},
		db:   db,
	}, nil
}

func (m *PgSQLRepository) WithLogger(h *slog.Logger) *PgSQLRepository {
	m.l = h
	return m
}

func (m *PgSQLRepository) StoreService(c context.Context, s Service) error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	_, err := m.db.ExecContext(
		c,
		"INSERT INTO service (ext_id, name, version, port, active) VALUES (?, ?, ?, ?, ?)",
		s.ext_id, s.name, s.version, s.port, s.active)
	return err
}
