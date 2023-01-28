package rd

import (
	"fmt"
	"time"
	"web/domain"

	"github.com/go-redis/redis"
)

const (
	address  = "redis"
	port     = ":6379"
	password = ""
)

type caсheStore struct {
	aTokenTTL time.Duration
	rTokenTTL time.Duration
	*redis.Client
}

func NewRedisClient(aToken, rToken time.Duration) (domain.CaсheStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address + port,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &caсheStore{
		aTokenTTL: aToken,
		rTokenTTL: rToken,
		Client:    client,
	}, nil
}

func (c *caсheStore) FindToken(id int64, token string) bool {
	key := fmt.Sprintf("user:%d", id)

	value, err := c.Get(key).Result()
	if err != nil {
		return false
	}

	return token == value
}

func (c *caсheStore) InsertToken(id int64, token string) error {
	key := fmt.Sprintf("user:%d", id)

	return c.Set(key, token, c.rTokenTTL).Err()
}

func (c *caсheStore) GetAccessTokenTTL() time.Duration {
	return c.aTokenTTL
}

func (c *caсheStore) GetRefreshTokenTTL() time.Duration {
	return c.rTokenTTL
}
