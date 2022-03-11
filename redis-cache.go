package main

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v7"
)

type redisCache struct {
	hosts    []string
	db      int
	expires time.Duration
	client  redis.ClusterClient
}

func NewRedisCache(hosts []string, db int, exp time.Duration) PostCache {
	return &redisCache{
		hosts:    hosts,
		db:      db,
		expires: exp,
		client:  *redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:     hosts,
		}),
	}
}


func (cache *redisCache) Set(key string, post string) {
	client := cache.client

	// serialize Post object to JSON
	json, err := json.Marshal(post)
	if err != nil {
		panic(err)
	}

	client.Set(key, json, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string) (string, error) {
	client := cache.client

	val, err := client.Get(key).Result()
	if err != nil {
		return "",nil
	}

	post := ""
	err = json.Unmarshal([]byte(val), &post)
	if err != nil {
		panic(err)
	}

	return post,nil
}
