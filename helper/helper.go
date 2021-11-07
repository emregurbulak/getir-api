package helper

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var Cache = cache.New(5*time.Minute, 5*time.Minute)

type InMemory []struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func SetCache(key string, memo interface{}) bool {
	Cache.Set(key, memo, cache.NoExpiration)
	return true
}

func GetCache(key string) (InMemory, bool) {
	var memo InMemory
	var found bool
	data, found := Cache.Get(key)

	if found {
		memo = data.(InMemory)
	}
	return memo, found
}
