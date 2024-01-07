package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	eb "github.com/dan-solli/homeapps/common/clients/eventbroker"
	"github.com/dan-solli/homeapps/microservice/servicemesh/config"
	"github.com/dan-solli/homeapps/microservice/servicemesh/service"
	pb "github.com/dan-solli/homeapps/proto/servicemesh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type server struct {
	pb.UnimplementedServiceMeshServiceServer
}

func NewGRpcServer(c config.GRpc) error {
	var opts []grpc.ServerOption

	if err := config.Viper().Unmarshal(&c); err != nil {
		log.Info("Failed to unmarshal config file", "err", err)
		panic(err)
	}
	log.Info("Hydrated config", "cfg", c)

	log.Info("Given port", "port", c.Grpc_port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", c.Grpc_port))
	if err != nil {
		log.Error("failed to listen", "err", err)
	}
	defer lis.Close()

	if c.Tls {
		creds, err := credentials.NewServerTLSFromFile(c.Certfile, c.Keyfile)
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

func (s *server) Announce(ctx context.Context, in *pb.AnnounceRequest) (*pb.AnnounceResponse, error) {
	tmpport, err := GetFreePort()
	if err != nil {
		log.Error("No free port number to hand out.", "err", err)
		return nil, err
	}

	sc := service.NewService(in.Name, in.Version, tmpport)

	if err := app.Store.StoreService(ctx, sc); err != nil {
		log.Error("Failed to save service to db", "err", err)
	}

	json := "[1, 2, 3]"

	ebc, err := eb.NewEventBrokerClient(app.cfg.Client.Eventbroker_address)
	if err != nil {
		log.Error("Failed to initialize eventbroker client", "err", err)
		return nil, err
	}

	r, err := ebc.AnnounceEvent(ctx, json)
	if err != nil {
		log.Error("Failed to announce event", "err", err)
	}
	log.Debug("PostEvent response:",
		"event_id", r.EventId,
		"corr_id", r.CorrelationId,
		"timestamp", r.CreatedAt.AsTime(),
	)

	return &pb.AnnounceResponse{
		Id:          sc.ExtId.String(),
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

///
///
///

func StartGRpcServer(wg *sync.WaitGroup, s *grpc.Server, c config.GRpc) *grpc.Server {
	log.Info("Given port", "port", c.Grpc_port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", c.Grpc_port))
	if err != nil {
		log.Error("failed to listen", "err", err)
	}
	defer lis.Close()
	pb.RegisterServiceMeshServiceServer(s, &server{})
	log.Info("server listening", "port", lis.Addr())

	if err := s.Serve(lis); err != grpc.ErrServerStopped {
		log.Error("failed to serve", "err", err)
	}
	return s
}

func CreateGRpcServer(c config.GRpc) (*grpc.Server, error) {
	var opts []grpc.ServerOption

	if err := config.Viper().Unmarshal(&c); err != nil {
		log.Info("Failed to unmarshal config file", "err", err)
		panic(err)
	}
	log.Info("Hydrated config", "cfg", c)

	if c.Tls {
		creds, err := credentials.NewServerTLSFromFile(c.Certfile, c.Keyfile)
		if err != nil {
			log.Error("Failed to generate credentials", "err", err)
			return nil, err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	s := grpc.NewServer(opts...)

	return s, nil
}
