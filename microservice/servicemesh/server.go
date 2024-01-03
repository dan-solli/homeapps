package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"log/slog"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	eb "github.com/dan-solli/homeapps/common/clients/eventbroker"
	"github.com/dan-solli/homeapps/microservice/servicemesh/service"
	pb "github.com/dan-solli/homeapps/proto/servicemesh"
)

// TODO: Need function to check health on registered services (and match with db)
// TODO: How and when should environment variables be read?
// TODO: Implement retries and waiting for remote services not responding.

var (
	log *slog.Logger
)

type ServiceMesh struct {
	cfg   ServiceMeshConfig
	gc    *eb.EventBrokerClient
	store service.IStore
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

	app.cfg = NewServiceMeshConfig()

	if err := NewGRpcServer(app.cfg); err != nil {
		log.Error("Failed to initialize grpc-server", "err", err)
	}

	ebc, err := eb.NewEventBrokerClient(viper.GetString("EVENTBROKER_ADDRESS"))
	if err != nil {
		log.Error("Failed to initialize eventbroker client", "err", err)
	}
	app.gc = ebc

	dbConfig := *service.NewDBConfig()

	if store, err := service.NewPgSQLRepository(dbConfig); err != nil {
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
	tmpport, err := GetFreePort()
	if err != nil {
		log.Error("No free port number to hand out.", "err", err)
		return nil, err
	}

	sc := service.NewService(in.Name, in.Version, tmpport)

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
		Id:          sc.GetExternalID(),
		Serviceport: int32(tmpport),
	}, nil
}

func (s *server) Denounce(ctx context.Context, in *pb.DenounceRequest) (*pb.DenounceResponse, error) {
	return &pb.DenounceResponse{Status: true}, nil
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
