package cachemap

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/kurnosovmak/short-link/pkg/logging"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Cache struct {
	Logger logging.Logger
	data   map[string]*CacheItem
	lock   sync.Mutex
}

type CacheContract interface {
	Get(key interface{}) interface{}
	Set(key, value interface{}, ttl time.Duration) error
	Clear(key interface{}) error
	hashKey(key interface{}) (string, error)
}
type CacheItem struct {
	Item    interface{}
	Expires time.Time
}

func New(logger logging.Logger) *Cache {
	return &Cache{
		Logger: logger,
		lock:   sync.Mutex{},
		data:   make(map[string]*CacheItem),
	}
}

func (cache *Cache) Get(key interface{}) interface{} {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	hash, err := cache.hashKey(key)

	if err != nil {
		cache.Logger.Logf(logrus.InfoLevel, "fn=get  key=%q status=error error=%q", hash, err)
		return nil
	}

	item := cache.data[hash]

	if item == nil {
		cache.Logger.Logf(logrus.InfoLevel, "fn=get key=%q status=miss", hash)
		return nil
	}

	if item.Expires.Before(time.Now()) {
		cache.Logger.Logf(logrus.InfoLevel, "fn=get key=%q status=expired", hash)
		_ = cache.Clear(hash)
		return nil
	}

	cache.Logger.Logf(logrus.InfoLevel, "fn=get key=%q status=hit", hash)
	return item.Item
}

func (cache *Cache) Set(key, value interface{}, ttl time.Duration) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	hash, err := cache.hashKey(key)

	if err != nil {
		cache.Logger.Logf(logrus.InfoLevel, "fn=set key=%q status=error error=%q", hash, err)
		return err
	}

	cache.data[hash] = &CacheItem{
		Item:    value,
		Expires: time.Now().Add(ttl),
	}

	cache.Logger.Logf(logrus.InfoLevel, "fn=set key=%q status=success", hash)
	return nil
}

func (cache *Cache) Clear(key interface{}) error {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	hash, err := cache.hashKey(key)

	if err != nil {
		cache.Logger.Logf(logrus.InfoLevel, "fn=clear key=%q status=error error=%q", hash, err)
		return err
	}

	if cache.data != nil && cache.data[hash] != nil {
		delete(cache.data, hash)
	}

	return nil
}

func (cache *Cache) hashKey(key interface{}) (string, error) {
	data, err := json.Marshal(key)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(data))[0:32], nil
}
