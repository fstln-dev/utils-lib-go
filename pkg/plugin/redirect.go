package plugin

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func Redirect() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil && se.Code >= 300 && se.Code < 400 {
				if tr, ok := transport.FromServerContext(ctx); ok {
					tr.ReplyHeader().Set("Location", se.Message)
				}
			}
			return
		}
	}
}

// func Redirect() middleware.Middleware {
// 	return func(handler middleware.Handler) middleware.Handler {
// 		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
// 			if tr, ok := transport.FromServerContext(ctx); ok {
// 				// Do something on entering
// 				defer func() {
// 					// Do something on exiting

// 					fmt.Println("redirect: ", tr)

// 					location := tr.ReplyHeader().Get("Location")
// 					if location != "" {
// 						err = errors.Newf(304, "", "")
// 						tr.ReplyHeader().Set("Location", location)
// 					}
// 				}()
// 			}
// 			reply, err = handler(ctx, req)

// 			if se := errors.FromError(err); er
// 		}
// 	}
// }
