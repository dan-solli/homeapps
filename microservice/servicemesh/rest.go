package main

import (
	"fmt"
	"net/http"

	"github.com/dan-solli/homeapps/microservice/servicemesh/config"
)

func NewHttpServer(c config.Http) error {
	if err := config.Viper().Unmarshal(&c); err != nil {
		log.Info("Failed to unmarshal config file", "err", err)
		panic(err)
	}
	log.Info("Hydrated config", "cfg", c)

	log.Info("Starting http-server")
	log.Info("Handler for metrics", "path", c.Metrics_path)
	http.HandleFunc(c.Metrics_path, getMetricsSpoof())
	log.Info("Handler for health", "path", c.Health_path)
	http.HandleFunc(c.Health_path, getHealth())

	err := http.ListenAndServe(fmt.Sprintf(":%d", c.Rest_port), nil)
	if err != nil {
		log.Error("failed to listen", "err", err)
	}
	return err
}

func getHealth() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func getMetricsSpoof() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
