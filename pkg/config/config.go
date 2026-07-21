package config

import (
	"context"
	"encoding/json"
	"errors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Option is etcd config option.
type Option func(o *options)

type options struct {
	ctx    context.Context
	path   string
	prefix bool
}

// Context with registry context.
func Context(ctx context.Context) Option {
	return Option(func(o *options) {
		o.ctx = ctx
	})
}

// Path is config path
func Path(p string) Option {
	return Option(func(o *options) {
		o.path = p
	})
}

// Prefix is config prefix
func Prefix(prefix bool) Option {
	return Option(func(o *options) {
		o.prefix = prefix
	})
}

type source struct {
	client  *clientv3.Client
	options *options
}

func New(c *clientv3.Client, opts ...Option) *source {
	options := &options{
		ctx:    context.Background(),
		path:   "",
		prefix: false,
	}

	for _, opt := range opts {
		opt(options)
	}

	return &source{
		client:  c,
		options: options,
	}
}

var ErrKvsNotFound = errors.New("kvs is Nil")

// Load return the config values
func (s *source) Load(m interface{}) error {
	var opts []clientv3.OpOption
	if s.options.prefix {
		opts = append(opts, clientv3.WithPrefix())
	}

	rsp, err := s.client.Get(s.options.ctx, s.options.path, opts...)
	if err != nil {
		return err
	}

	if len(rsp.Kvs) < 1 {
		return ErrKvsNotFound
	}

	return json.Unmarshal(rsp.Kvs[0].Value, m)
}

func (s *source) Close() error {
	return s.client.Close()
}

// Watch return the watcher
func (s *source) Watch(m interface{}) error {

	if err := s.Load(m); err != nil {
		return err
	}

	go func(s *source, m interface{}) {

		w := newWatcher(s)
		for {
			w.Next(m)
		}

	}(s, m)
	return nil
}
