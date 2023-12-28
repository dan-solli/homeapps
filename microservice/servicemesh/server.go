package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"database/sql"
	"log/slog"

	pb "github.com/dan-solli/homeapps/microservice/servicemesh/public"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/google/uuid"
)

// TODO: Need function to check health on registered services (and match with db)

type serviceCache struct {
	ext_id  uuid.UUID
	name    string
	version string
	port    int32
	active  bool
}

type server struct {
	pb.UnimplementedServiceMeshServiceServer
}

type runtimeConfig struct {
	db       *sql.DB
	tls      bool
	certFile string
	keyFile  string
	port     int
}

var (
	rtc runtimeConfig
	svc []serviceCache
	log *slog.Logger
)

func init() {
	viper.SetEnvPrefix("MS_SM")
	viper.AutomaticEnv()

	rtc = runtimeConfig{}
	svc = []serviceCache{}

	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	init_flags(&rtc)

	// Init database
	err := init_db(&rtc)
	if err != nil {
		log.Error("Fatal error database", "err", err)
	}
	defer rtc.db.Close()

	num, err := readState()
	if err != nil {
		log.Error("Failed reading state from database", "err", err)
	}
	log.Info("Read services from database", "cnt", num)

	// Init prometheus metrics
	// Init tracer
}

func main() {
	lis, err := init_server(6000)
	if err != nil {
		log.Error("failed to start server", "err", err)
	}

	opts := init_serveropts(&rtc)

	s := grpc.NewServer(opts...)
	pb.RegisterServiceMeshServiceServer(s, &server{})
	log.Info("server listening", "port", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Error("failed to serve", "err", err)
	}
}

func init_flags(rtc *runtimeConfig) {
	rtc.tls = viper.GetBool("TLS")
	rtc.certFile = viper.GetString("CERTFILE")
	rtc.keyFile = viper.GetString("KEYFILE")
	rtc.port = viper.GetInt("GRPC_PORT")
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

func init_serveropts(rtc *runtimeConfig) []grpc.ServerOption {
	if rtc.tls {
		creds, err := credentials.NewServerTLSFromFile(rtc.certFile, rtc.keyFile)
		if err != nil {
			log.Error("Failed to generate credentials", "err", err)
		}
		return []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		return []grpc.ServerOption{}
	}
}

func init_server(port int) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Error("failed to listen", "err", err)
	}
	return lis, err
}

func readState() (int, error) {
	rows, err := rtc.db.Query("SELECT id, name, version, port FROM service WHERE active = true")

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

func (s *server) Announce(ctx context.Context, in *pb.AnnounceRequest) (*pb.AnnounceResponse, error) {
	sc := serviceCache{
		ext_id:  uuid.New(),
		name:    in.GetName(),
		version: in.GetVersion(),
		active:  true,
	}

	rows := rtc.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(port), ?) FROM service WHERE active = true",
		viper.GetInt("SERVICE_PORT_RANGE_START"))

	var tmpport int32

	err := rows.Scan(&tmpport)
	if err != nil {
		log.Error("Failed to get free port number from database", "err", err)
	}
	sc.port = int32(tmpport + 1)

	result, err := rtc.db.ExecContext(
		ctx,
		"INSERT INTO service (ext_id, name, version, port, active) VALUES (?, ?, ?, ?, ?)",
		sc.ext_id, sc.name, sc.version, sc.port, sc.active)
	if err != nil {
		log.Error("Failed to save service to db", "service", sc.name, "port", sc.port, "err", err)
	}
	if touchedRows, err := result.RowsAffected(); touchedRows == 0 {
		log.Error("Failed to save service to db", "service", sc.name, "port", sc.port, "err", err)
	}

	svc = append(svc, sc)

	/*
		TODO: Post Event about the newcomer.
	*/
	return &pb.AnnounceResponse{
		Id:          sc.ext_id.String(),
		Serviceport: sc.port,
	}, nil
}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{
		Id:        "1",
		ServiceId: "1",
		Port:      6001,
	}, nil
}

func (s *server) Deregister(ctx context.Context, in *pb.DeregisterRequest) (*pb.DeregisterResponse, error) {
	return &pb.DeregisterResponse{Status: true}, nil
}

func (s *server) Denounce(ctx context.Context, in *pb.DenounceRequest) (*pb.DenounceResponse, error) {
	return &pb.DenounceResponse{Status: true}, nil
}
