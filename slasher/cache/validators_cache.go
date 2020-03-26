package cache

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// validatorsCacheSize defines the max number of validators public keys the cache can hold.
	validatorsCacheSize = 300000
	// Metrics for the validator cache.
	validatorsCacheHit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "validators_cache_hit",
		Help: "The total number of cache hits on the validators cache.",
	})
	validatorsCacheMiss = promauto.NewCounter(prometheus.CounterOpts{
		Name: "validators_cache_miss",
		Help: "The total number of cache misses on the validators cache.",
	})
)

// ValidatorsCache is used to store the public keys needed for signature verification.
type ValidatorsCache struct {
	cache *lru.Cache
}

// NewValidatorsCache initializes the cache.
func NewValidatorsCache(size int, onEvicted func(key interface{}, value interface{})) (*ValidatorsCache, error) {
	if size != 0 {
		validatorsCacheSize = size
	}
	cache, err := lru.NewWithEvict(validatorsCacheSize, onEvicted)
	if err != nil {
		return nil, err
	}
	return &ValidatorsCache{cache: cache}, nil
}

// Get returns an ok bool and the cached value for the requested validator id key, if any.
func (c *ValidatorsCache) Get(validatorIdx uint64) ([]byte, bool) {
	item, exists := c.cache.Get(validatorIdx)
	if exists && item != nil {
		validatorsCacheHit.Inc()
		return item.([]byte), true
	}

	validatorsCacheMiss.Inc()
	return nil, false
}

// Set the response in the cache.
func (c *ValidatorsCache) Set(validatorIdx uint64, publicKey []byte) {
	_ = c.cache.Add(validatorIdx, publicKey)
}

// Delete removes a validator id from the cache and returns if it existed or not.
// Performs the onEviction function before removal.
func (c *ValidatorsCache) Delete(validatorIdx uint64) bool {
	return c.cache.Remove(validatorIdx)
}

// Has returns true if the key exists in the cache.
func (c *ValidatorsCache) Has(validatorIdx uint64) bool {
	return c.cache.Contains(validatorIdx)
}

// Clear removes all keys from the ValidatorCache.
func (c *ValidatorsCache) Clear() {
	c.cache.Purge()
}
