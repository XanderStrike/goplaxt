package store

import (
	"github.com/xanderstrike/goplaxt/lib/user"
)

// Store is the interface for All the store types
type Store interface {
	WriteUser(user user.User)
	GetUser(id string) user.User
}

// Utils
func flatTransform(s string) []string { return []string{} }
