package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gookit/slog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	autenRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "homeapps_ms_users_application_auten_requests_total",
		Help: "The total number of authentication requests",
	})
)

type AutenticationService interface {
	Autenticate(ctx context.Context, username, password string) (string, error)
}

type autenticationService struct{}

func (autenticationService) Autenticate(ctx context.Context, username, password string) (string, error) {
	return "", nil
}

func main() {
	slog.Info("hello world")
	fmt.Println("hello world")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

	autenRequests.Inc()
}
