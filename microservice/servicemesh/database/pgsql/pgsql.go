package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type PgSQLRepository struct {
	lock *sync.RWMutex
	db   *sql.DB
	l    *slog.Logger
}

type Service struct {
	Ext_id  uuid.UUID
	Name    string
	Version string
	Port    int32
	Active  bool
}

func NewPgSQLRepository(l *slog.Logger) (*PgSQLRepository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetInt("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASS"),
		viper.GetString("DB_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		l.Error("Fatal error getting database connection", "err", err)
		return nil, err
	}

	return &PgSQLRepository{
		lock: &sync.RWMutex{},
		db:   db,
		l:    l,
	}, nil
}

func (m *PgSQLRepository) GetFreePortNumber(c context.Context) (int32, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	rows := m.db.QueryRowContext(c, "SELECT COALESCE(MAX(port), ?) FROM service WHERE active = true",
		viper.GetInt("SERVICE_PORT_RANGE_START"))

	var tmpport int32

	err := rows.Scan(&tmpport)
	if err != nil {
		m.l.Error("Failed to get free port number from database", "err", err)
		return -1, err
	}
	return int32(tmpport + 1), nil
}

func (m *PgSQLRepository) StoreService(c context.Context, s Service) error {
	_, err := m.db.ExecContext(
		c,
		"INSERT INTO service (ext_id, name, version, port, active) VALUES (?, ?, ?, ?, ?)",
		s.Ext_id, s.Name, s.Version, s.Port, s.Active)
	return err
}
