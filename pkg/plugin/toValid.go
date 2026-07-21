package plugin

import (
	"context"
	"time"

	"github.com/fstln-dev/utils-lib-go/pkg/jwt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

/**
* 接口验证插件
* API 使用 JWT
*   入口：验证 Header
*
*
* GRPC 预定使用 JWT
*   入口：暂无，预定验证 metadata,MD
*
*   // Kratos 已处理出口 GRPC 的 MD 写入
*   出口：暂无，预定写入 metadata.MD 中
**/

// to Valid
func ToValid(secret string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				// Do something on entering
				hk := tr.RequestHeader().Get("x-fstln-token")

				if hk == "" {
					return nil, errors.Unauthorized("NOT_FOUND_AUTH", "Must have authorized")
				}

				// 鉴权
				ok, err := jwt.JWTValid(secret, hk)
				if !ok {
					if err != nil {
						// log
					}
					return false, errors.InternalServer("FAIL_VERIFICATE", "Verification failure")
				}

				m := &JWTTenantModel{}
				err = jwt.JWTDecode(hk, m)
				if err != nil {
					// log
					return false, errors.InternalServer("FAIL_DECODE", "Verification failure")
				}

				if m.Iss != "fstln" || m.Exp.IsZero() {
					return false, errors.InternalServer("FAIL_VERIFICATE", "Verification failure")
				}

				if m.Exp.Before(time.Now()) {
					return false, errors.InternalServer("TIME_OUT", "Login has timed out, please log in again")
				}

				return handler(ctx, req)
			}

			return nil, errors.InternalServer("FAIL_VERIFICATE", "Not find the http context")
		}
	}
}

func SetGrpcValid() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				// Do something on entering
				tr.RequestHeader().Set("JWT", "写入 Jwt")

				defer func() {
					// Do something on exiting
				}()
			}
			return handler(ctx, req)
		}
	}
}
