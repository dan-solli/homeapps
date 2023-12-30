package main

import (
	"context"
	"sync"

	_ "github.com/lib/pq"
)

type PgSQLServiceMeshRepository struct {
	lock        *sync.RWMutex
	serviceMesh services.ServiceMesh
}

func NewPgSQLServiceMeshRepository(sm services.ServiceMesh) *PgSQLServiceMeshRepository {
	if sm.IsZero() {
		log.Error("missing smFactory")
	}

	return &PgSQLServiceMeshRepository{
		lock:        &sync.RWMutex{},
		serviceMesh: sm,
	}
}

func (m PgSQLServiceMeshRepository) GetFreePortNumber(_ context.Context) (int32, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.serviceMesh.getFreePortNumber()
}

func (m PgSQLServiceMeshRepository) ReadState() (int32, error) {
	return m.serviceMesh.readState(), nil
}
