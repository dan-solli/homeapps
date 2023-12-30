package main

import "database/sql"

type runtimeConfig struct {
	db       *sql.DB
	tls      bool
	certFile string
	keyFile  string
	port     int
}

var (
	rtc runtimeConfig
)
