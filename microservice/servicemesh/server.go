package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"sync"

	"log/slog"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	eb "github.com/dan-solli/homeapps/common/clients/eventbroker"
	pb "github.com/dan-solli/homeapps/proto/servicemesh"
)

// TODO: Need function to check health on registered services (and match with db)
// TODO: How and when should environment variables be read?
// TODO: Implement retries and waiting for remote services not responding.

var (
	log *slog.Logger
)

type PgSQLRepository struct {
	lock *sync.RWMutex
	db   *sql.DB
}

type iStore interface {
	getFreePortNumber(c context.Context) (int32, error)
	storeService(c context.Context, s service) error
}

type ServiceMesh struct {
	cfg   ServiceMeshConfig
	gc    *eb.EventBrokerClient
	store iStore
}

type server struct {
	pb.UnimplementedServiceMeshServiceServer
}

type service struct {
	ext_id  uuid.UUID
	name    string
	version string
	port    int32
	active  bool
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

	if err := NewGRpcServer(app.cfg); err != nil {
		log.Error("Failed to initialize grpc-server", "err", err)
	}

	ebc, err := eb.NewEventBrokerClient(viper.GetString("EVENTBROKER_ADDRESS"))
	if err != nil {
		log.Error("Failed to initialize eventbroker client", "err", err)
	}
	app.gc = ebc

	if store, err := NewPgSQLRepository(app.cfg); err != nil {
		log.Error("Failed to initialize repository", "err", err)
	} else {
		app.store = store
	}

	// Init prometheus metrics
	// Init tracer
}

func NewGRpcServer(c ServiceMeshConfig) error {
	var opts []grpc.ServerOption

	lis, err := NewServerListener(c.grpcport)
	if err != nil {
		log.Error("failed to start server", "err", err)
	}

	if c.tls {
		creds, err := credentials.NewServerTLSFromFile(c.certFile, c.keyFile)
		if err != nil {
			log.Error("Failed to generate credentials", "err", err)
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	s := grpc.NewServer(opts...)
	pb.RegisterServiceMeshServiceServer(s, &server{})
	log.Info("server listening", "port", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Error("failed to serve", "err", err)
		return err
	}
	return nil
}

func NewServerListener(port int) (net.Listener, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Error("failed to listen", "err", err)
		return nil, err
	}
	return lis, nil
}

func NewPgSQLRepository(c ServiceMeshConfig) (*PgSQLRepository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetInt("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASS"),
		viper.GetString("DB_NAME"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error("Fatal error getting database connection", "err", err)
		return nil, err
	}

	return &PgSQLRepository{
		lock: &sync.RWMutex{},
		db:   db,
	}, nil
}

func (m PgSQLRepository) getFreePortNumber(c context.Context) (int32, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	rows := m.db.QueryRowContext(c, "SELECT COALESCE(MAX(port), ?) FROM service WHERE active = true",
		viper.GetInt("SERVICE_PORT_RANGE_START"))

	var tmpport int32

	err := rows.Scan(&tmpport)
	if err != nil {
		log.Error("Failed to get free port number from database", "err", err)
		return -1, err
	}
	return int32(tmpport + 1), nil
}

func (m PgSQLRepository) storeService(c context.Context, s service) error {
	_, err := m.db.ExecContext(
		c,
		"INSERT INTO service (ext_id, name, version, port, active) VALUES (?, ?, ?, ?, ?)",
		s.ext_id, s.name, s.version, s.port, s.active)
	return err
}

func (s *server) Announce(ctx context.Context, in *pb.AnnounceRequest) (*pb.AnnounceResponse, error) {
	tmpport, err := app.store.getFreePortNumber(ctx)
	if err != nil {
		log.Error("No free port number to hand out.", "err", err)
		return nil, err
	}

	sc := service{
		ext_id:  uuid.New(),
		name:    in.Name,
		version: in.Version,
		port:    tmpport,
		active:  true,
	}

	if err := app.store.storeService(ctx, sc); err != nil {
		log.Error("Failed to save service to db", "err", err)
	}

	json := "[1, 2, 3]"

	r, err := app.gc.AnnounceEvent(ctx, json)
	if err != nil {
		log.Error("Failed to announce event", "err", err)
	}
	log.Debug("PostEvent response:",
		"event_id", r.EventId,
		"corr_id", r.CorrelationId,
		"timestamp", r.CreatedAt.AsTime(),
	)

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
