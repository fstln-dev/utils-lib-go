package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fstln-dev/utils-lib-go/pkg/graphql/jsonutil"
	"github.com/fstln-dev/utils-lib-go/pkg/http"
)

const (
	OperationQuery        Operation = "query"
	OperationMutation     Operation = "mutation"
	OperationSubscription Operation = "subscription"
)

type Operation string
type Variables map[string]interface{}

type GraphQLData struct {
	Query     string    `json:"query"`
	Variables Variables `json:"variables,omitempty"`
}

type errors []struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
	Path []string
}

func (e errors) Error() string {
	return fmt.Sprintf("%s: %s", strings.Join(e[0].Path, "."), e[0].Message)
}

type graphOut struct {
	Data       *json.RawMessage `json:"data"`
	Errors     errors           `json:"errors"`
	Extensions interface{}      `json:"extensions"`
}

func Header(h map[string]string) http.Option {
	return http.Header(h)
}

func Request(ctx context.Context, path string, op Operation, query interface{}, variables Variables, opts ...http.Option) error {

	// response 结构
	resp := &graphOut{}

	// 拼接 request, 发送 post 请求
	if err := http.HttpsPost(
		ctx, path, resp,
		append(opts, http.Request(&GraphQLData{
			Query:     constructQuery(op, query, variables),
			Variables: variables,
		}))...,
	); err != nil {
		return err
	}

	if resp.Data != nil {
		if err := jsonutil.UnmarshalGraphQL(*resp.Data, query); err != nil {
			return err
		}
	}
	if len(resp.Errors) > 0 {
		return resp.Errors
	}
	return nil
}
