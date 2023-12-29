package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"database/sql"
	"log/slog"

	eb "github.com/dan-solli/homeapps/common/proto/eventbroker"
	pb "github.com/dan-solli/homeapps/common/proto/servicemesh"

	"github.com/golang/protobuf/jsonpb"
	structpb "github.com/golang/protobuf/ptypes/struct"

	"google.golang.org/protobuf/types/known/timestamppb"
	//"google.golang.org/protobuf/types/known/struct"

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
	ebc eb.EventBrokerServiceClient
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
	var g_opts []grpc.DialOption
	conn, err := grpc.Dial(viper.GetString("EVENTBROKER_ADDRESS"), g_opts...)
	if err != nil {
		log.Error("Can't connect to EventBroker", "err", err)
	}
	defer conn.Close()
	ebc = eb.NewEventBrokerServiceClient(conn)

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

func json2pb(json string) (*eb.EventPayload, error) {
	msg := &eb.EventPayload{Data: &structpb.Value{}}
	jsm := jsonpb.Unmarshaler{}

	if err := jsm.Unmarshal(bytes.NewReader([]byte(json)), msg.Data); err != nil {
		return nil, err
	}
	return msg, nil
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

	json := "[1, 2, 3]"
	pbjson, err := json2pb(json)
	if err != nil {
		log.Error("Could not convert string json to object", "err", err, "json", json)
	}

	r, err := ebc.PostEvent(context.Background(), &eb.PostEventRequest{
		EventId:       uuid.New().String(),
		CorrelationId: uuid.New().String(),
		Source:        "Microservice:ServiceMesh",
		Event:         "framework.service.announce",
		CreatedAt:     timestamppb.New(time.Now()),
		Payload: &eb.EventPayload{
			ContentType: "text/json",
			Data:        pbjson.Data,
		},
	})
	if err != nil {
		log.Error("Call to PostEvent failed.", "request", "<bleh>", "err", err)
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
