package backend

import "github.com/andresoro/hemera/pkg/cache"

// Backend interface handles flushing data from cache to a receiver i.e graphite
// Do not worry about clearing the cache from the backend interface as it is done at the server layer
// this allows multiple backends to be used during a purge cycle
type Backend interface {
	Purge(c *cache.Cache) error
}
