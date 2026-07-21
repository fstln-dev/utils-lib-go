package grpc

import (
	"context"
	"time"

	"github.com/go-kratos/etcd/registry"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/selector/filter"
	kratos_grpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

/**
* Data 服务发现模块
* 通过 GetServiceGRPCClient({服务名称}, {etcd 服务注册客户端}) 获取
* 使用: 保存在 xxxRepo 中, 子方法调用
*
* Data Layer, service discovery module
* Get by GetServiceGRPCClient({service name}, {etcd service registration client})
* Usage: Saved in xxxRepo, submethod call
**/

type ClientConn struct {
	Client *grpc.ClientConn
}

func GetServiceGRPCClient(rr *registry.Registry, service, version string, timeout time.Duration, log log.Logger, m ...middleware.Middleware) (*ClientConn, error) {

	gClient, err := kratos_grpc.DialInsecure(
		context.Background(),
		kratos_grpc.WithDiscovery(rr),
		kratos_grpc.WithEndpoint("discovery:///"+service),
		kratos_grpc.WithTimeout(timeout),
		kratos_grpc.WithNodeFilter(filter.Version(version)),
		kratos_grpc.WithLogger(log),
		kratos_grpc.WithMiddleware(m...),
	)
	if err != nil {
		return nil, err
	}

	return &ClientConn{
		Client: gClient,
	}, nil
}

func (c ClientConn) Close() error {
	return c.Client.Close()
}
