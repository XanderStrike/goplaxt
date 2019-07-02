package store

import (
	"context"
)

// Store is the interface for All the store types
type Store interface {
	WriteUser(user User)
	GetUser(id string) *User
	Ping(ctx context.Context) error
}

// Utils
func flatTransform(s string) []string { return []string{} }
