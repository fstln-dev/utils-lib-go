package etcd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

// Option is etcd config option.
type Option func(o *options)

type options struct {
	dialTimeout time.Duration
	block       bool
}

func DialTimeout(t time.Duration) Option {
	return Option(func(o *options) {
		o.dialTimeout = t
	})
}

func Block(b bool) Option {
	return Option(func(o *options) {
		o.block = b
	})
}

func New(addr []string, opts ...Option) (*clientv3.Client, error) {
	options := &options{
		dialTimeout: time.Second * 10,
		block:       true,
	}

	for _, opt := range opts {
		opt(options)
	}

	g := []grpc.DialOption{}
	if options.block {
		g = append(g, grpc.WithBlock())
	}

	client, err := clientv3.New(
		clientv3.Config{
			Endpoints:   addr,
			DialTimeout: options.dialTimeout,
			DialOptions: g,
		})

	if err != nil {
		return nil, err
	}

	return client, nil
}
