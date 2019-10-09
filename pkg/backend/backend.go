package backend

import "github.com/andresoro/hemera/pkg/cache"

// Backend interface handles flushing data from cache to a reciever i.e graphite
type Backend interface {
	Purge(c *cache.Cache)
}
