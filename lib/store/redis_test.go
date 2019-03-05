package store

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := NewRedisClient(s.Addr())
	store := NewRedisStore(client)

	s.HSet("user:id123", "username", "halkeye")
	s.HSet("user:id123", "access", "access123")
	s.HSet("user:id123", "refresh", "refresh123")
	s.HSet("user:id123", "updated", "02-25-2019")

	expected, err := json.Marshal(&User{
		ID:           "id123",
		Username:     "halkeye",
		AccessToken:  "access123",
		RefreshToken: "refresh123",
		Updated:      time.Date(2019, 02, 25, 0, 0, 0, 0, time.UTC),
	})
	actual, err := json.Marshal(store.GetUser("id123"))

	assert.EqualValues(t, string(expected), string(actual))
}
