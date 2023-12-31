package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/dan-solli/homeapps/microservice/servicemesh/config"
)

type PgSQLRepository struct {
	lock *sync.RWMutex
	db   *sql.DB
	l    *slog.Logger
}

func NewPgSQLRepository(d config.DB, log *slog.Logger) (*PgSQLRepository, error) {

	if err := config.Viper().Unmarshal(&d); err != nil {
		log.Info("Failed to unmarshal config file", "err", err)
		panic(err)
	}
	log.Info("Hydrated config", "cfg", d)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Db_host, d.Db_port, d.Db_user, d.Db_pass, d.Db_name)

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
		s.ExtId, s.Name, s.Version, s.Port, s.Active)
	return err
}

func (m *PgSQLRepository) GetServices(c context.Context) ([]Service, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	rows, err := m.db.QueryContext(c, "SELECT * FROM service")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var s Service
		err := rows.Scan(&s.ExtId, &s.Name, &s.Version, &s.Port, &s.Active)
		if err != nil {
			return nil, err
		}
		services = append(services, s)
	}
	return services, nil
}
