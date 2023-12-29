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

/*
type RuntimeConfig struct {
	db       *sql.DB
	tls      bool
	certFile string
	keyFile  string
	port     int
}
*/
/*
type ServiceCache struct {
	ext_id  uuid.UUID
	name    string
	version string
	port    int32
	active  bool
}
*/

var (
	//rtc RuntimeConfig
	//svc []serviceCache
	log *slog.Logger
)

func init() {
	viper.SetEnvPrefix("MS_SM")
	viper.AutomaticEnv()

	init_config()
	init_cache()

	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	init_flags(&rtc)

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
	rtc.tls = viper.GetBool("TLS")
	rtc.certFile = viper.GetString("CERTFILE")
	rtc.keyFile = viper.GetString("KEYFILE")
	rtc.port = viper.GetInt("GRPC_PORT")
}

func init_server(port int) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Error("failed to listen", "err", err)
	}
	return lis, err
}
