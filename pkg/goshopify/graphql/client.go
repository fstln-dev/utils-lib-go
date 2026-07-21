package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
	"regexp"
	"strconv"
	"time"
)

const (
	// todo 如何设置http超时时间
	defaultHttpTimeout = 10
)

var ErrRecordNotFound = errors.New("record not found")

type Option func(*Client)

type Client struct {
	ShopName   string
	Token      string
	ApiVersion string

	// todo 重试次数
}

// PageInfo 分页信息
type PageInfo struct {
	StartCursor     string `json:"startCursor"`
	EndCursor       string `json:"endCursor"`
	HasNextPage     bool   `json:"hasNextPage"`
	HasPreviousPage bool   `json:"hasPreviousPage"`
}

// UserErrors Mutation操作的错误信息
type UserErrors []struct {
	Field   []graphql.String
	Message graphql.String
}

type MoneyV2 struct {
	Amount       string `json:"amount"`       //todo Decimal 类型
	CurrencyCode string `json:"currencyCode"` //todo CurrencyCode 类型
}

type Job struct {
	Id   string `json:"id"`
	Done bool   `json:"done"`
}

func (c *Client) Path() string {
	return fmt.Sprintf("https://%s.myshopify.com/admin/api/%s/graphql.json", c.ShopName, c.ApiVersion)
}

func (c *Client) Query(ctx context.Context, query interface{}, variables graphql.Variables) error {
	return graphql.Request(ctx, c.Path(), graphql.OperationQuery, query, variables,
		graphql.Header(map[string]string{"X-Shopify-Access-Token": c.Token}))
}

func (c *Client) Mutation(ctx context.Context, query interface{}, variables graphql.Variables) error {
	return graphql.Request(ctx, c.Path(), graphql.OperationMutation, query, variables,
		graphql.Header(map[string]string{"X-Shopify-Access-Token": c.Token}))
}

func NewClient(shopName string, token string, opts ...Option) (*Client, error) {
	conf := &Client{
		ShopName:   shopName,
		Token:      token,
		ApiVersion: getQuarterlyVersion(),
	}

	for _, opt := range opts {
		opt(conf)
	}

	return conf, nil
}

func getQuarterlyVersion() string {
	year, month, _ := time.Now().Date()
	switch month {
	case 2, 3, 4:
		return fmt.Sprintf("%d-01", year)
	case 5, 6, 7:
		return fmt.Sprintf("%d-04", year)
	case 8, 9, 10:
		return fmt.Sprintf("%d-07", year)
	case 11, 12:
		return fmt.Sprintf("%d-10", year)
	default: // month == 1
		return fmt.Sprintf("%d-10", year-1)
	}
}

func Ids2Gids(format string, ids []int64) []string {
	var gids []string
	for i := range ids {
		gids = append(gids, fmt.Sprintf(format, ids[i]))
	}
	return gids
}

func IdByGraphQL(graphqlId string) int64 {
	matches := regexp.MustCompile(`(\d+)$`).FindStringSubmatch(graphqlId)
	if len(matches) > 1 {
		id, err := strconv.ParseInt(matches[1], 10, 64)
		if err == nil {
			return id
		}
	}
	return 0
}

func IdStrByGraphQL(graphqlId string) string {
	matches := regexp.MustCompile(`([^/]+)$`).FindStringSubmatch(graphqlId)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func PtrFloat(val float64) *float64 {
	return &val
}
func PtrInt(val int) *int {
	return &val
}
func PtrBool(val bool) *bool {
	return &val
}
func PtrString(val string) *string {
	return &val
}
func PtrTime(val time.Time) *time.Time {
	return &val
}
