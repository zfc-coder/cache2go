/*
 * Simple caching library with expiration capabilities
 *     Copyright (c) 2012, Radu Ioan Fericean
 *                   2013-2017, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package cache2go

import (
	"sync"
)

var (
	cache = make(map[string]any)
	mutex sync.RWMutex
)

type Key comparable
type Value any

// Cache returns the existing cache table with given name or creates a new one
// if the table does not exist yet.
func Cache[K Key, V Value](table string) *CacheTable[K, V] {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		mutex.Lock()
		t, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			t = &CacheTable[K, V]{
				name:  table,
				items: make(map[K]*CacheItem[K, V]),
			}
			cache[table] = t
		}
		mutex.Unlock()
	}

	return t.(*CacheTable[K, V])
}
