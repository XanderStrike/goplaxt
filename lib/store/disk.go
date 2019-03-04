package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/peterbourgon/diskv"
	"github.com/xanderstrike/goplaxt/lib/user"
)

// DiskStore is a storage engine that writes to the disk
type DiskStore struct{}

// WriteUser will write a user object to disk
func (s DiskStore) WriteUser(user user.User) {
	s.writeField(user.ID, "username", user.Username)
	s.writeField(user.ID, "access", user.AccessToken)
	s.writeField(user.ID, "refresh", user.RefreshToken)
	s.writeField(user.ID, "updated", user.Updated.Format("01-02-2006"))
}

// GetUser will load a user from disk
func (s DiskStore) GetUser(id string) user.User {
	updated, _ := time.Parse("01-02-2006", s.readField(id, "updated"))
	user := user.User{
		ID:           id,
		Username:     strings.ToLower(s.readField(id, "username")),
		AccessToken:  s.readField(id, "access"),
		RefreshToken: s.readField(id, "refresh"),
		Updated:      updated,
	}

	return user
}

func (s DiskStore) writeField(id, field, value string) {
	err := s.write(fmt.Sprintf("%s.%s", id, field), value)
	if err != nil {
		panic(err)
	}
}

func (s DiskStore) readField(id, field string) string {
	value, err := s.read(fmt.Sprintf("%s.%s", id, field))
	if err != nil {
		panic(err)
	}
	return value
}

func (s DiskStore) write(key, value string) error {
	d := diskv.New(diskv.Options{
		BasePath:     "keystore",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	return d.Write(key, []byte(value))
}

func (s DiskStore) read(key string) (string, error) {
	d := diskv.New(diskv.Options{
		BasePath:     "keystore",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	value, err := d.Read(key)
	return string(value), err
}
