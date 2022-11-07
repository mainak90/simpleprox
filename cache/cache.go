package cache

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"sync"
)

var NotCachable = errors.New("This object is not cachable. Skipping!")

type CacheEntry struct {
	res Result
	ready chan struct{}
	invalid bool
}

type Result struct {
	value interface{}
	err   error
}

type Memo struct {
	requestFunc Func
	cache 		map[string]*CacheEntry
	mutex 		sync.Mutex
}

type Func func(key string) (interface{}, error)

func NewCache(f Func) *Memo {
	return &Memo{requestFunc: f, cache: make(map[string]*CacheEntry)}
}

func (memo *Memo) GetKey(key string) (interface{}, error) {
	memo.mutex.Lock()
	cachedkey := memo.cache[key]
	if cachedkey == nil {
		log.WithFields(log.Fields{"Controller": "cache"}).Info("Key not found ", key)
		cachedkey = &CacheEntry{ready: make(chan struct{})}
		memo.cache[key] = cachedkey
		memo.mutex.Unlock()
		log.WithFields(log.Fields{"Controller": "cache"}).Info("Cache key not found", key)

		tres, terr := memo.requestFunc(key)
		if terr != nil {
			if terr == NotCachable {
				cachedkey.invalid = true
				//we see from header that this url is not suitable for cache
				//so we do not update this response to corresponding cache entry
				return tres, nil
			} else {
				//error when requesting
				//todo: consider retrying
				cachedkey.res.value, cachedkey.res.err = nil, terr
			}
		} else {
			//save response to cache entry
			cachedkey.res.value, cachedkey.res.err = tres, nil
		}
		close(cachedkey.ready)
	} else if cachedkey.invalid {
		memo.mutex.Unlock()
		//this url is not suitable for cache, we know this from the real response
		return nil, NotCachable
	} else {
		memo.mutex.Unlock()
		<-cachedkey.ready
	}
	return cachedkey.res.value, cachedkey.res.err
}