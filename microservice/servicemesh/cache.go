package main

import (
	"github.com/google/uuid"
)

type serviceCache struct {
	ext_id  uuid.UUID
	name    string
	version string
	port    int32
	active  bool
}
