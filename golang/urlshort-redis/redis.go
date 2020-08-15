package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"crypto/sha1"

	"github.com/go-redis/redis/v8"
	"github.com/mattheath/base62"
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

// Shorten convert url to shortlink
func (r RedisCli) Shorten(url string, exp int64) (string, error) {
	// convert url to sha1 hash
	h := toSha1(url)

	// fetch it if the url is cached
	d, err := r.Cli.Get(context.Background(), fmt.Sprintf(URLHashKey, h)).Result()
	if err == redis.Nil {
		// key not existed, nothing to do
	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
			// expiration, nothing to do
		} else {
			return d, nil
		}
	}

	// increase the global counter
	err = r.Cli.Incr(context.Background(), URLIDKEY).Err()
	if err != nil {
		return "", err
	}

	// encode global counter to base62
	id, err := r.Cli.Get(context.Background(), URLIDKEY).Int64()
	if err != nil {
		return "", nil
	}

	eid := base62.EncodeInt64(id)

	// store the irl against this encoded id
	err = r.Cli.Set(context.Background(), fmt.Sprintf(ShortlinkKey, eid), url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", nil
	}

	// store the url against the hash if it
	err = r.Cli.Set(context.Background(), fmt.Sprintf(URLHashKey, h), eid, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", nil
	}

	detail, err := json.Marshal(&URLDetail{
		URL:                 url,
		CreatedAt:           time.Now().String(),
		ExpirationInMinutes: time.Duration(exp),
	})

	if err != nil {
		return "", err
	}

	// stor the url detail against this encoded id
	err = r.Cli.Set(context.Background(), fmt.Sprintf(ShortlinkDetailKey, eid), detail,
		time.Minute*time.Duration(exp)).Err()

	if err != nil {
		return "", nil
	}

	return eid, nil
}

func toSha1(str string) string {
	return string(sha1.New().Sum([]byte(str)))
}

// ShortlinkInfo returns the details of the shortlink
func (r RedisCli) ShortlinkInfo(eid string) (interface{}, error) {
	d, err := r.Cli.Get(context.Background(), fmt.Sprintf(ShortlinkDetailKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{400, errors.New("Unknown short URL")}
	} else if err != nil {
		return nil, err
	} else {
		return d, nil
	}
}

// Unshorten convert shortlink to url
func (r RedisCli) Unshorten(eid string) (string, error) {
	url, err := r.Cli.Get(context.Background(), fmt.Sprintf(ShortlinkKey, eid)).Result()
	if err == redis.Nil {
		return "", StatusError{404, err}
	} else if err != nil {
		return "", err
	} else {
		return url, nil
	}
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
