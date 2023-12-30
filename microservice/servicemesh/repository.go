package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Repo struct {
	Db *sql.DB
}

func RepoInterface(db *sql.DB) *Repo {
	return &Repo{Db: db}
}

func (repo *Repo) Connect() (bool, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("DB_HOST"), viper.GetInt("DB_PORT"), viper.GetString("DB_USER"), viper.GetString("DB_PASS"), viper.GetString("DB_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Fatal error database", "err", err)
	}
	rtc.db = db
	return true, nil
}

func SqlConfig() *sql.DB {
	var Db *sql.DB

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("DB_HOST"), viper.GetInt("DB_PORT"), viper.GetString("DB_USER"), viper.GetString("DB_PASS"), viper.GetString("DB_NAME"))
	Db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Could not estalibhs connection to database", "err", err)
		return nil
	}
	return Db

}

func init_db2() *sql.DB {
	d := sqlx.database
	db := sql.new
}

func init_db(rtc *runtimeConfig) error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("DB_HOST"), viper.GetInt("DB_PORT"), viper.GetString("DB_USER"), viper.GetString("DB_PASS"), viper.GetString("DB_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Fatal error database", "err", err)
	}
	rtc.db = db
	return err
}

func readState(d *sql.DB) (int, error) {
	rows, err := d.Query("SELECT id, name, version, port FROM service WHERE active = true")

	if err != nil {
		log.Error("Failed to fetch state from database", "err", err)
		return -1, err
	}
	defer rows.Close()

	var sv serviceCache
	counter := 0

	for rows.Next() {
		if err := rows.Scan(&sv.ext_id, &sv.name, &sv.version, &sv.port); err != nil {
			log.Error("Failed to get row of data", "err", err)
		}
		svc = append(svc, sv)
		counter++
	}
	return counter, nil
}

func getFreePort(c context.Context, d *sql.DB) (int32, error) {
	rows := d.QueryRowContext(c, "SELECT COALESCE(MAX(port), ?) FROM service WHERE active = true",
		viper.GetInt("SERVICE_PORT_RANGE_START"))

	var tmpport int32

	err := rows.Scan(&tmpport)
	if err != nil {
		log.Error("Failed to get free port number from database", "err", err)
		return -1, err
	}
	return int32(tmpport + 1), nil
}

func storeService(c context.Context, d *sql.DB, s serviceCache) error {
	_, err := d.ExecContext(
		c,
		"INSERT INTO service (ext_id, name, version, port, active) VALUES (?, ?, ?, ?, ?)",
		s.ext_id, s.name, s.version, s.port, s.active)
	return err
}
