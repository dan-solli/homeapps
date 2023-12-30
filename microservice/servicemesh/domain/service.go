package service

import (
	"github.com/google/uuid"
)

type MeshService struct {
	ext_id  uuid.UUID
	name    string
	version string
	port    int32
	active  bool
}

func NewServiceMeshCache() ([]MeshService, error) {
	return []MeshService{}, nil
}

/*
type meshCache struct {
	services []serviceCache
}

func (m meshCache) findServiceByName(name string) (*serviceCache, error) {
	for _, s := range m.services {
		if s.name == name {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("Service %q not found", name)
}
*/
