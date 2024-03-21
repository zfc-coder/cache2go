/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"sync"
	"time"
)

// CacheItem is an individual cache item
// Parameter data contains the user-set value in the cache.
type CacheItem[K Key, V Value] struct {
	sync.RWMutex

	// The item's key.
	key K
	// The item's data.
	data V
	// How long will the item live in the cache when not being accessed/kept alive.
	lifeSpan time.Duration

	// Creation timestamp.
	createdOn time.Time
	// Last access timestamp.
	accessedOn time.Time
	// How often the item was accessed.
	accessCount int64

	// Callback method triggered right before removing the item from the cache
	aboutToExpire []func(key K)
}

// NewCacheItem returns a newly created CacheItem.
// Parameter key is the item's cache-key.
// Parameter lifeSpan determines after which time period without an access the item
// will get removed from the cache.
// Parameter data is the item's value.
func NewCacheItem[K Key, V Value](key K, lifeSpan time.Duration, data V) *CacheItem[K, V] {
	t := time.Now()
	return &CacheItem[K, V]{
		key:           key,
		lifeSpan:      lifeSpan,
		createdOn:     t,
		accessedOn:    t,
		accessCount:   0,
		aboutToExpire: nil,
		data:          data,
	}
}

// KeepAlive marks an item to be kept for another expireDuration period.
func (item *CacheItem[K, V]) KeepAlive() {
	item.Lock()
	defer item.Unlock()
	item.accessedOn = time.Now()
	item.accessCount++
}

// LifeSpan returns this item's expiration duration.
func (item *CacheItem[K, V]) LifeSpan() time.Duration {
	// immutable
	return item.lifeSpan
}

// AccessedOn returns when this item was last accessed.
func (item *CacheItem[K, V]) AccessedOn() time.Time {
	item.RLock()
	defer item.RUnlock()
	return item.accessedOn
}

// CreatedOn returns when this item was added to the cache.
func (item *CacheItem[K, V]) CreatedOn() time.Time {
	// immutable
	return item.createdOn
}

// AccessCount returns how often this item has been accessed.
func (item *CacheItem[K, V]) AccessCount() int64 {
	item.RLock()
	defer item.RUnlock()
	return item.accessCount
}

// Key returns the key of this cached item.
func (item *CacheItem[K, V]) Key() K {
	// immutable
	return item.key
}

// Data returns the value of this cached item.
func (item *CacheItem[K, V]) Data() V {
	// immutable
	return item.data
}

// SetAboutToExpireCallback configures a callback, which will be called right
// before the item is about to be removed from the cache.
func (item *CacheItem[K, V]) SetAboutToExpireCallback(f func(K)) {
	if len(item.aboutToExpire) > 0 {
		item.RemoveAboutToExpireCallback()
	}
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = append(item.aboutToExpire, f)
}

// AddAboutToExpireCallback appends a new callback to the AboutToExpire queue
func (item *CacheItem[K, V]) AddAboutToExpireCallback(f func(K)) {
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = append(item.aboutToExpire, f)
}

// RemoveAboutToExpireCallback empties the about to expire callback queue
func (item *CacheItem[K, V]) RemoveAboutToExpireCallback() {
	item.Lock()
	defer item.Unlock()
	item.aboutToExpire = nil
}
