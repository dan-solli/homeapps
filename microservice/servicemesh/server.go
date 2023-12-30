package main

import (
	"fmt"
	"net"
	"os"

	"log/slog"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// TODO: Need function to check health on registered services (and match with db)
// TODO: How and when should environment variables be read?
// TODO: Implement retries and waiting for remote services not responding.

var (
	log *slog.Logger
)

func init() {
	viper.SetEnvPrefix("MS_SM")
	viper.AutomaticEnv()

	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	rtc.tls = viper.GetBool("TLS")
	rtc.certFile = viper.GetString("CERTFILE")
	rtc.keyFile = viper.GetString("KEYFILE")
	rtc.port = viper.GetInt("GRPC_PORT")

	// TODO: Is there a reason these two functions has a different return signature?
	if _, err := init_grpc_client(); err != nil {
		log.Error("Failed to initialize grpc-client", "err", err)
	}
	if err := init_grpc_server(&rtc); err != nil {
		log.Error("Failed to initialize grpc-server", "err", err)
	}

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

	// Init prometheus metrics
	// Init tracer
}

func main() {
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
