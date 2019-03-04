package store

import (
	"testing"

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
	user := store.GetUser("halkeye")

	assert.Equal(t, user, nil, "thingie")
}
