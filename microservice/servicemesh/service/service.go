package service

import "github.com/google/uuid"

type Service struct {
	ext_id  uuid.UUID
	name    string
	version string
	port    int
	active  bool
}

func NewService(name string, version string, port int) Service {
	return Service{
		ext_id:  uuid.New(),
		name:    name,
		version: version,
		port:    port,
		active:  true,
	}
}

func (s *Service) GetExternalID() string {
	return s.ext_id.String()
}
