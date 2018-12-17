package store

import (
	"fmt"
	"os"
	"time"

	"github.com/peterbourgon/diskv"
)

type User struct {
	ID           string
	Username     string
	AccessToken  string
	RefreshToken string
	Updated      time.Time
}

func NewUser(username, accessToken, refreshToken string) User {
	id := uuid()
	user := User{
		ID:           id,
		Username:     username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Updated:      time.Now(),
	}
	writeField(id, "username", username)
	writeField(id, "access", accessToken)
	writeField(id, "refresh", refreshToken)
	writeField(id, "updated", user.Updated.Format("01-02-2006"))
	return user
}

func writeField(id, field, value string) {
	err := write(fmt.Sprintf("%s.%s", id, field), value)
	if err != nil {
		panic(err)
	}
}

func GetUser(id string) User {
	updated, _ := time.Parse("01-02-2006", readField(id, "updated"))
	user := User{
		ID:           id,
		Username:     readField(id, "username"),
		AccessToken:  readField(id, "access"),
		RefreshToken: readField(id, "refresh"),
		Updated:      updated,
	}

	return user
}

func readField(id, field string) string {
	value, err := read(fmt.Sprintf("%s.%s", id, field))
	if err != nil {
		panic(err)
	}
	return value
}

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
