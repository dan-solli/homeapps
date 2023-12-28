package common

import (
	uuid "github.com/satori/go.uuid"
)

// UUID generates a new UUID
func UUID() string {
	return uuid.NewV4().String()
}
