// Package multicache provides a "fallback" cache implementation that
// short-circuits gets and writes/deletes to all underlying caches.
package multicache

// Fallback is a cache that wraps a list of caches. Sets write to the first
// WaitNWrites caches. Gets read from the caches in sequence until a cache entry
// is found. Deletes delete from all caches.
type Fallback struct {
	caches      []Underlying
	WaitNWrites int
}

func (f *Fallback) Get(key string) (resp []byte, ok bool) {
	for _, c := range f.caches {
		resp, ok = c.Get(key)
		if ok {
			return
		}
	}
	return
}

func (f *Fallback) Set(key string, resp []byte) {
	for i, c := range f.caches {
		if i < f.WaitNWrites {
			c.Set(key, resp)
		} else {
			go c.Set(key, resp)
		}
	}
}

func (f *Fallback) Delete(key string) {
	for _, c := range f.caches {
		c.Delete(key)
	}
}

// NewFallback returns a new Fallback cache with WaitNWrites == len(caches).
func NewFallback(caches ...Underlying) *Fallback {
	return &Fallback{caches: caches, WaitNWrites: len(caches)}
}
