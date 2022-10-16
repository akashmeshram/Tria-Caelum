package main


import (
	"context"
	"os"
	"runtime"
	"time"

	"github.com/allegro/bigcache/v3"
	cacheStoreCache "github.com/eko/gocache/cache"
	cacheStore "github.com/eko/gocache/store"
	"github.com/patrickmn/go-cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Cache caching
type Cache interface {
	Flush(context.Context) error
	Get(context.Context, string, interface{}) error
	Set(context.Context, string, interface{}) error
	SetTTL(context.Context, time.Duration, string, interface{}) error
	Delete(context.Context, string) error
}


//BigCache ...
func BigCache(eviction time.Duration) (Cache, error) {
	bigcacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(eviction))
	if err != nil {
		return nil, err
	}
	bigcacheStore := cacheStore.NewBigcache(bigcacheClient, nil)

	cacheManager := cacheStoreCache.New(bigcacheStore)
	return &LocalBigCache{
		db: cacheManager,
	}, nil
}

func Local() Cache {
	gocacheClient := cache.New(5*time.Minute, 10*time.Minute)
	gocacheStore := cacheStore.NewGoCache(gocacheClient, nil)

	cacheManager := cacheStoreCache.New(gocacheStore)
	return &localCache{
		db: *cacheManager,
	}
}