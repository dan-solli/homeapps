package eventbroker

import (
	"context"
	"log/slog"
	"os"
	"time"

	eb "github.com/dan-solli/homeapps/common/proto/eventbroker"
	spb "github.com/golang/protobuf/ptypes/struct"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	log = slog.New(slog.NewJSONHandler(os.Stdout, nil))
)

type EventBrokerClient struct {
	c eb.EventBrokerServiceClient
}

func NewEventBrokerClient(adr string) (*EventBrokerClient, error) {
	opts := []grpc.DialOption{}

	conn, err := grpc.Dial(adr, opts...)
	if err != nil {
		log.Error("Can't connect to EventBroker", "err", err)
		return nil, err
	}
	defer conn.Close()
	return &EventBrokerClient{
		c: eb.NewEventBrokerServiceClient(conn),
	}, nil
}

func (ebc EventBrokerClient) announceEvent(c context.Context, j string) (*eb.PostEventResponse, error) {
	field, err := json2pb(j)
	if err != nil {
		log.Error("Could not convert string json to object", "err", err, "json", j)
		return nil, err
	}

	r, err := ebc.c.PostEvent(context.Background(), &eb.PostEventRequest{
		EventId:       uuid.New().String(),
		CorrelationId: uuid.New().String(),
		Source:        "Microservice:ServiceMesh",
		Event:         "framework.service.announce",
		CreatedAt:     timestamppb.New(time.Now()),
		Payload: &eb.EventPayload{
			ContentType: "text/json",
			Data:        field.Data,
		},
	})
	if err != nil {
		log.Error("Call to PostEvent failed.", "request", "<bleh>", "err", err)
		return nil, err
	}
	log.Debug("Successful post event call", "event_id", r.EventId, "corrid", r.CorrelationId)
	return r, err
}

func json2pb(json string) (*eb.EventPayload, error) {
	msg := &eb.EventPayload{Data: &spb.Value{}}

	if err := protojson.Unmarshal([]byte(json), msg.Data); err != nil {
		return nil, err
	}
	return msg, nil
}
