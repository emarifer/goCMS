package api

import (
	"fmt"
	"sync/atomic"
	"time"

	shardedmap "github.com/zutto/shardedmap"
)

// Do not store over this amount
// of MBs in the cache
const MAX_CACHE_SIZE_MB = 10

type EndpointCache struct {
	Name       string
	Contents   []byte
	ValidUntil time.Time
}

func emptyEndpointCache() EndpointCache {

	return EndpointCache{"", []byte{}, time.Now()}
}

type TimeValidator struct{}

type CacheValidator interface {
	IsValid(cache *EndpointCache) bool
}

func (v *TimeValidator) IsValid(cache *EndpointCache) bool {

	// We only return the cache if it's still valid
	return cache.ValidUntil.After(time.Now())
}

type TimedCache struct {
	CacheMap      shardedmap.ShardMap
	CacheTimeout  time.Duration
	EstimatedSize atomic.Uint64 // in bytes
	Validator     CacheValidator
}

type Cache interface {
	Get(name string) (EndpointCache, error)
	Store(name string, buffer []byte) error
	Size() uint64
}

func MakeCache(n_shards int, expiry_duration time.Duration) Cache {

	return &TimedCache{
		CacheMap:      shardedmap.NewShardMap(n_shards),
		CacheTimeout:  expiry_duration,
		EstimatedSize: atomic.Uint64{},
		Validator:     &TimeValidator{},
	}
}

// Implementing the interface Cache

func (cache *TimedCache) Get(name string) (EndpointCache, error) {
	cacheEntry := cache.CacheMap.Get(name)
	// if the endpoint is cached
	if cacheEntry != nil {
		cacheContents, ok := (*cacheEntry).(EndpointCache)
		if !ok {
			return emptyEndpointCache(), fmt.Errorf(
				"it has not been possible to make a type assertion on the cache stored under that key",
			)
		}

		// We only return the cache if it's still valid
		if cache.Validator.IsValid(&cacheContents) {
			return cacheContents, nil
		} else {
			cache.CacheMap.Delete(name)

			return emptyEndpointCache(), fmt.Errorf(
				"cached endpoint had expired",
			)
		}
	}

	return emptyEndpointCache(), fmt.Errorf("cache does not contain key")
}

func (cache *TimedCache) Store(name string, buffer []byte) error {
	// Only store to the cache if we have enough space left
	afterSizeMB := float64(
		cache.EstimatedSize.Load()+uint64(len(buffer)),
	) / 1000000
	if afterSizeMB > MAX_CACHE_SIZE_MB {
		return fmt.Errorf("maximum size reached")
	}

	// Saving to shardedmap
	var cacheEntry interface{} = EndpointCache{
		Name:       name,
		Contents:   buffer,
		ValidUntil: time.Now().Add(cache.CacheTimeout),
	}

	cache.CacheMap.Set(name, &cacheEntry)
	cache.EstimatedSize.Add(uint64(len(buffer)))

	return nil
}

func (cache *TimedCache) Size() uint64 {

	return cache.EstimatedSize.Load()
}
