package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

var (
	GlobalTtl    string = "45m"
	FavoritesKey string = "stonk-favorites:%v"
	ErrNotFound  error  = errors.New("Favorites Not Found")
)

func InitRedisClient(connAddr, pass string) (*redis.Client, error) {

	dbOpts := &redis.Options{
		Addr:     connAddr,
		Password: pass,
		DB:       0,
	}

	client := redis.NewClient(dbOpts)
	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to Redis. host: %s, error: %v", connAddr, err)
	}

	return client, nil
}

func GetFavorites(ctx context.Context, c *redis.Client, userId string) ([]string, error) {

	userKey := fmt.Sprintf(FavoritesKey, userId)
	vals, err := c.SMembers(userKey).Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, ErrNotFound
		}
		return []string{}, err
	}

	return vals, nil
}

func AddFavorite(ctx context.Context, c *redis.Client, userId, symbol string) error {
	userKey := fmt.Sprintf(FavoritesKey, userId)

	_, err := c.SAdd(userKey, symbol).Result()
	return err
}

func RemoveFavorite(ctx context.Context, c *redis.Client, userId, symbol string) error {
	userKey := fmt.Sprintf(FavoritesKey, userId)

	_, err := c.SRem(userKey, symbol).Result()
	return err
}
