package main

import (
	"fmt"
	"net"
	"os"

	"log/slog"

	service "github.com/dan-solli/homeapps/microservice/servicemesh/domain"

	_ "github.com/lib/pq"
)

// TODO: Need function to check health on registered services (and match with db)
// TODO: How and when should environment variables be read?
// TODO: Implement retries and waiting for remote services not responding.

var (
	log *slog.Logger
)

type ServiceMesh struct {
	cfg      ServiceMeshConfig
	gc       EventBrokerClient
	services []service.MeshService
}

var (
	app = ServiceMesh{}
)

func main() {
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var err error

	app.cfg, err = NewServiceMeshConfig()
	if err != nil {
		log.Error("Failed to set up config", "err", err)
	}
	// TODO: Is there a reason these two functions has a different return signature?
	app.gc, err = NewGRPCClient(app.cfg)

	if err := init_grpc_server(app.cfg); err != nil {
		log.Error("Failed to initialize grpc-server", "err", err)
	}

	/*
		// Init database
		err := init_db(&rtc)
		if err != nil {
			log.Error("Fatal error database", "err", err)
		}
		defer rtc.db.Close()

		num, err := readState(rtc.db)
		if err != nil {
			log.Error("Failed reading state from database", "err", err)
		}
		log.Info("Read services from database", "cnt", num)
	*/
	// Init prometheus metrics
	// Init tracer
}

func init_flags(rtc *runtimeConfig) {

}

func init_server(port int) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Error("failed to listen", "err", err)
	}
	return lis, err
}
