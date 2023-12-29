package main

import (
	"github.com/google/uuid"
)

var svc []serviceCache

type serviceCache struct {
	ext_id  uuid.UUID
	name    string
	version string
	port    int32
	active  bool
}

func init_cache() {
	svc = []serviceCache{}
}
