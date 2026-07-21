package alert

import (
	"context"
	"fmt"

	"github.com/fstln-dev/utils-lib-go/pkg/http"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

var ErrUnknownRequest = errors.InternalServer("UNKNOWN", "unknown request error")

type robotResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type robotRequest struct {
	Msgtype  string        `json:"msgtype"`
	Markdown robotMarkdown `json:"markdown"`
}
type robotMarkdown struct {
	Content string `json:"content"`
}

func Alert(server_name string, alert_link string, ignore map[string]interface{}) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code   int32
				reason string
				// kind      string
				// operation string
			)

			// if info, ok := transport.FromServerContext(ctx); ok {
			// 	kind = info.Kind().String()
			// 	operation = info.Operation()
			// }
			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Reason
			}

			if _, ok := ignore[reason]; ok {
				return
			}

			if len(alert_link) > 0 && code >= 400 {
				go func() {
					http.HttpsPost(
						context.TODO(),
						alert_link,
						&robotResponse{},
						http.Request(
							&robotRequest{
								Msgtype: "markdown",
								Markdown: robotMarkdown{
									Content: fmt.Sprintf("<font color=\"warning\">%s 接口报错</font>\n ERROR CODE: %d\n ERROR REASON: %s\n ERROR DETAIL: %s", server_name, code, reason, err.Error()),
								},
							},
						),
					)
				}()
			}

			return
		}
	}
}

func AlertPanicHandle(server_name string, alert_link string) func(ctx context.Context, req, err interface{}) error {
	return func(ctx context.Context, req, err interface{}) error {

		go func() {
			http.HttpsPost(
				context.TODO(),
				alert_link,
				&robotResponse{},
				http.Request(
					&robotRequest{
						Msgtype: "markdown",
						Markdown: robotMarkdown{
							Content: fmt.Sprintf("<font color=\"warning\">%s 接口报错</font>\n PANIC ERROR: %v: %+v", server_name, err, req),
						},
					},
				),
			)
		}()

		return ErrUnknownRequest
	}
}
