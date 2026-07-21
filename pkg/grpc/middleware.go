package grpc

import (
	"context"

	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
)

func AddServiceNameToMetadata(service_name string) middleware.Middleware {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			ctx = metadata.AppendToClientContext(ctx, "X-FSTLN-SERVICE-NAME", service_name)
			return h(ctx, req)
		}
	}
}
