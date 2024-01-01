package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"log/slog"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	eb "github.com/dan-solli/homeapps/common/clients/eventbroker"
	pg "github.com/dan-solli/homeapps/microservice/servicemesh/database/pgsql"
	pb "github.com/dan-solli/homeapps/proto/servicemesh"
)

// TODO: Need function to check health on registered services (and match with db)
// TODO: How and when should environment variables be read?
// TODO: Implement retries and waiting for remote services not responding.

var (
	log *slog.Logger
)

type iStore interface {
	GetFreePortNumber(c context.Context) (int32, error)
	StoreService(c context.Context, s pg.Service) error
}

type ServiceMesh struct {
	cfg   ServiceMeshConfig
	gc    *eb.EventBrokerClient
	store iStore
}

type server struct {
	pb.UnimplementedServiceMeshServiceServer
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

	if store, err := pg.NewPgSQLRepository(log); err != nil {
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

func (s *server) Announce(ctx context.Context, in *pb.AnnounceRequest) (*pb.AnnounceResponse, error) {
	tmpport, err := app.store.GetFreePortNumber(ctx)
	if err != nil {
		log.Error("No free port number to hand out.", "err", err)
		return nil, err
	}

	sc := pg.Service{
		Ext_id:  uuid.New(),
		Name:    in.Name,
		Version: in.Version,
		Port:    tmpport,
		Active:  true,
	}

	if err := app.store.StoreService(ctx, sc); err != nil {
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
		Id:          sc.Ext_id.String(),
		Serviceport: sc.Port,
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
