package data

import (
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type Option func(*options)
type options struct {
	pw           string
	db           int
	writeTimeout time.Duration
	readTimeout  time.Duration
}

func Password(password string) Option {
	return Option(func(o *options) {
		o.pw = password
	})
}

func Db(db int) Option {
	return Option(func(o *options) {
		o.db = db
	})
}

func WriteTimeout(time time.Duration) Option {
	return Option(func(o *options) {
		o.writeTimeout = time
	})
}

func ReadTimeout(time time.Duration) Option {
	return Option(func(o *options) {
		o.readTimeout = time
	})
}

// NewRedisClient is a func to New Client for Redis V8
func NewRedisClient(addr string, opts ...Option) (*redis.Client, error) {

	option := &options{
		db:           0,
		writeTimeout: time.Second,
		readTimeout:  time.Second,
	}

	for _, opt := range opts {
		opt(option)
	}

	if addr == "" {
		return nil, errors.New("Redis Addr Can not Nil")
	}

	c := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     option.pw,
		DB:           option.db,
		WriteTimeout: option.writeTimeout,
		ReadTimeout:  option.readTimeout,
	})
	// redis.AddHook(redisotel.TracingHook{})

	_, err := c.Ping(c.Context()).Result()
	if err != nil {
		return nil, err
	}

	return c, nil
}
