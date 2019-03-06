package store

import ()

// Store is the interface for All the store types
type Store interface {
	WriteUser(user User)
	GetUser(id string) *User
	Ping() error
}

// Utils
func flatTransform(s string) []string { return []string{} }
