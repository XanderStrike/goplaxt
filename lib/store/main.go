package store

import (
	"fmt"
	"os"

	"github.com/peterbourgon/diskv"
)

func flatTransform(s string) []string { return []string{} }

func write(key, value string) error {
	d := diskv.New(diskv.Options{
		BasePath:     "keystore",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	return d.Write(key, []byte(value))
}

func read(key string) (string, error) {
	d := diskv.New(diskv.Options{
		BasePath:     "keystore",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})
	value, err := d.Read(key)
	return string(value), err
}

func uuid() string {
	f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}

func NewUser(username, access_token, refresh_token string) (string, error) {
	id := uuid()
	err := write(fmt.Sprintf("%s.username", id), username)
	if err != nil {
		return id, err
	}
	err = write(fmt.Sprintf("%s.access", id), access_token)
	if err != nil {
		return id, err
	}
	err = write(fmt.Sprintf("%s.refresh", id), refresh_token)
	if err != nil {
		return id, err
	}
	return id, nil
}

func GetUser(id string) (string, string, string, error) {
	username, err := read(fmt.Sprintf("%s.username", id))
	if err != nil {
		return "", "", "", err
	}
	access_token, err := read(fmt.Sprintf("%s.access", id))
	if err != nil {
		return "", "", "", err
	}
	refresh_token, err := read(fmt.Sprintf("%s.refresh", id))
	if err != nil {
		return "", "", "", err
	}
	return username, access_token, refresh_token, nil
}
