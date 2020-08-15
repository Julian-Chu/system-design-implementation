package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// URLIDKEY is global counter
	URLIDKEY = "next.url.id"
	// ShortlinkKey mapping the shortlink to the url
	ShortlinkKey = "shortlink:%s:url"
	// URLHashKey mapping the hash of the url to the shortlink
	URLHashKey = "urlhash:%s:url"
	// ShortlinkDetailKey mapping the shortlink to the deital of url
	ShortlinkDetailKey = "shortlink:%s:detail"
)

// RedisCli contains a redis client
type RedisCli struct {
	Cli *redis.Client
}

// URLDetail contains the detail fo the shortlink
type URLDetail struct {
	URL                 string        `json:"url"`
	CreatedAt           string        `json:"created_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

// NewRedisCli create a redis Clitn
func NewRedicCli(addr string, passwd string, db int) *RedisCli {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	if _, err := c.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	return &RedisCli{Cli: c}
}
