package main

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/dan-solli/homeapps/common/proto/servicemesh"
)

type server struct {
	pb.UnimplementedServiceMeshServiceServer
}

func init_grpc_server(r *RuntimeConfig) error {
	var opts = []grpc.ServerOption{}

	lis, err := init_server(6000)
	if err != nil {
		log.Error("failed to start server", "err", err)
	}

	if r.tls {
		creds, err := credentials.NewServerTLSFromFile(r.certFile, r.keyFile)
		if err != nil {
			log.Error("Failed to generate credentials", "err", err)
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	} else {
		opts = []grpc.ServerOption{}
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
	sc := ServiceCache{
		ext_id:  uuid.New(),
		name:    in.GetName(),
		version: in.GetVersion(),
		active:  true,
	}

	tmpport, err := getFreePort(ctx, getConfig())
	if err != nil {
		log.Error("No free port number to hand out.", "err", err)
	}
	sc.port = tmpport

	if err := storeService(ctx, rtc, sc); err != nil {
		log.Error("Failed to save service to db", "service", sc.name, "port", sc.port, "err", err)
	}

	svc = append(svc, sc)

	json := "[1, 2, 3]"

	r, err := announceEvent(ctx, json)
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
