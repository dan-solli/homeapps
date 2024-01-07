package service

import "github.com/google/uuid"

type Service struct {
	ExtId   uuid.UUID `json:"ext_id"`
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Port    int       `json:"port"`
	Active  bool      `json:"-"`
}

/*
Name string `json:"full_name"`
Age int `json:"age,omitempty"`
Active bool `json:"-"`
lastLoginAt string
*/

func NewService(name string, version string, port int) Service {
	return Service{
		ExtId:   uuid.New(),
		Name:    name,
		Version: version,
		Port:    port,
		Active:  true,
	}
}
