package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"time"

	"github.com/allegro/bigcache"
)

type big_cache struct {
	cache *bigcache.BigCache
}

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string, value interface{}) (bool, error)
	Del(key string) error
}

type Option func(c *bigcache.Config)

func Shards(n int) Option {
	return func(c *bigcache.Config) {
		if n < 0 {
			return
		}
		c.Shards = n
	}
}

func TTL(t time.Duration) Option {
	return func(c *bigcache.Config) {
		c.LifeWindow = t
	}
}

func CleanTime(t time.Duration) Option {
	return func(c *bigcache.Config) {
		c.CleanWindow = t
	}
}

func Logger(l bigcache.Logger) Option {
	return func(c *bigcache.Config) {
		c.Logger = l
	}
}

func NewCache(opts ...Option) (Cache, error) {

	config := bigcache.Config{
		Shards:             16,
		LifeWindow:         10 * time.Second,
		CleanWindow:        5 * time.Second,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		HardMaxCacheSize:   500,
		Verbose:            false,
	}

	for _, f := range opts {
		f(&config)
	}

	ca, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}

	return &big_cache{
		cache: ca,
	}, nil
}

func (c *big_cache) Set(key string, value interface{}) error {

	bytes, err := serializeGOB(value)
	if err != nil {
		return err
	}

	return c.cache.Set(key, bytes)
}

func (c *big_cache) Get(key string, value interface{}) (bool, error) {

	bytes, err := c.cache.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return false, nil
		}
		return false, err
	}

	err = deserializeGOB(bytes, value)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (c *big_cache) Del(key string) error {
	err := c.cache.Delete(key)
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return err
	}
	return nil
}

// gob encode interface (point / struct)
func serializeGOB(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)

	if err := enc.Encode(value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// gob decode interface (point)
func deserializeGOB(valueBytes []byte, value interface{}) error {
	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(value); err != nil {
		return err
	}
	return nil
}
