package store

import (
	"fmt"
	"strings"
	"time"

	"github.com/xanderstrike/goplaxt/lib/user"

	"github.com/go-redis/redis"
)

// RedisStore is a storage engine that writes to redis
type RedisStore struct {
	client redis.Client
}

// NewRedisClient creates a new redis client object
func NewRedisClient(addr string) redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return *client
}

// NewRedisStore creates new store
func NewRedisStore(client redis.Client) RedisStore {
	return RedisStore{
		client: client,
	}
}

// WriteUser will write a user object to redis
func (s RedisStore) WriteUser(user user.User) {
	data := make(map[string]interface{})
	data["username"] = user.Username
	data["access"] = user.AccessToken
	data["refresh"] = user.RefreshToken
	data["updated"] = user.Updated.Format("01-02-2006")
	s.client.HMSet("user:"+user.ID, data)
}

// GetUser will load a user from redis
func (s RedisStore) GetUser(id string) user.User {
	data, _ := s.client.HGetAll("user:" + id).Result()
	fmt.Printf("Data: %v", data)
	updated, _ := time.Parse("01-02-2006", data["updated"])
	user := user.User{
		ID:           id,
		Username:     strings.ToLower(data["username"]),
		AccessToken:  data["access"],
		RefreshToken: data["refresh"],
		Updated:      updated,
	}

	return user
}
