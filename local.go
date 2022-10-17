package main

import (
	"context"
	"encoding/json"
	"fmt"
	gocache "github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	//DefaultExpiration ...
	DefaultExpiration time.Duration = 0
)

type LocalBigCache struct {
	db *gocache.Cache
}

type localCache struct {
	db gocache.Cache
}

//Flush ...
func (c *LocalBigCache) Flush(ctx context.Context) error {
	return c.db.Clear()
}

//Get ...
func (c *LocalBigCache) Get(ctx context.Context, key string, value interface{}) error {
	if v, err := c.db.Get(key); err == nil {
		if data, ok := v.([]byte); ok {
			if err := json.Unmarshal(data, value); err != nil {
				return err
			}
			return nil
		}
	}
	return ErrNotFound
}

//Set ...
func (c *LocalBigCache) Set(ctx context.Context, key string, value interface{}) error {
	return c.SetTTL(ctx, DefaultExpiration, key, value)
}



func (c *localCache) Flush(ctx context.Context) error {
	c.db.Clear()
	return nil
}

func (c *localCache) Get(ctx context.Context, key string, value interface{}) error {
	var ok error
	value, ok = c.db.Get(key)
	if ok != nil {
		fmt.Println("ok not nil: ", ok)
	}
	return nil
}

func (c *localCache) Set(ctx context.Context, key string, value interface{}) error {
	return c.SetTTL(ctx, cache.DefaultExpiration, key, value)
}
