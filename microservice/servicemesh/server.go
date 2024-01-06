package main

import (
	"os"

	"log/slog"

	_ "github.com/lib/pq"

	"github.com/dan-solli/homeapps/microservice/servicemesh/config"
	"github.com/dan-solli/homeapps/microservice/servicemesh/service"
)

// TODO: Need function to check health on registered services (and match with db)
// TODO: How and when should environment variables be read?
// TODO: Implement retries and waiting for remote services not responding.

type ServiceMesh struct {
	cfg   *config.Config
	Store service.IStore
}

var (
	log *slog.Logger
	app = ServiceMesh{}
)

func main() {
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// slog.LogLoggerLevel = slog.LevelDebug

	app.cfg = config.NewConfig(log)

	log.Info("Starting service mesh")

	go func() {
		log.Info("Starting grpc-server")
		if err := NewGRpcServer(app.cfg.Server.GRpc); err != nil {
			log.Error("Failed to initialize grpc-server", "err", err)
		}
	}()

	log.Info("Initializing repository")
	if store, err := service.NewPgSQLRepository(app.cfg.Server.DB, log); err != nil {
		log.Error("Failed to initialize repository", "err", err)
	} else {
		app.Store = store
	}

	if err := NewHttpServer(app.cfg.Server.Http); err != nil {
		log.Error("Failed to initialize http-server", "err", err)
	}

	// Init prometheus metrics
	// Init tracer
	log.Info("Service mesh started")
}
