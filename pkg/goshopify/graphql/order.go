package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/fstln-dev/utils-lib-go/pkg/goshopify/input"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
	"strings"
	"time"
)

type OrderRepo struct {
	client *Client
}

func NewOrderRepo(config *Client) *OrderRepo {
	return &OrderRepo{client: config}
}

type MoneyBag struct {
	ShopMoney        MoneyV2 `json:"shopMoney"`
	PresentmentMoney MoneyV2 `json:"presentmentMoney"`
}

type DiscountApplicationAllocationMethod string
type DiscountApplicationTargetSelection string
type DiscountApplicationTargetType string

type DiscountApplication struct {
	Index            int                                 `json:"index"`
	AllocationMethod DiscountApplicationAllocationMethod `json:"allocationMethod"`
	TargetSelection  DiscountApplicationTargetSelection  `json:"targetSelection"`
	TargetType       DiscountApplicationTargetType       `json:"targetType"`
}

type DiscountAllocation struct {
	AllocatedAmountSet  MoneyBag            `json:"allocatedAmountSet"`
	DiscountApplication DiscountApplication `json:"discountApplication"`
}
type LineItem struct {
	Id                  string               `json:"id"`
	Name                string               `json:"name"`
	Title               string               `json:"title"`
	Sku                 string               `json:"sku"`
	Quantity            int                  `json:"quantity"`
	TotalDiscountSet    MoneyBag             `json:"totalDiscountSet"`
	Variant             Variant              `json:"variant"`
	Product             Product              `json:"product"`
	DiscountAllocations []DiscountAllocation `json:"discountAllocations"`
}

type Refund struct {
	Id string `json:"id"`
}

type Fulfillment struct {
	DisplayStatus string    `json:"displayStatus"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeliveredAt   time.Time `json:"deliveredAt"`
}

type Order struct {
	CreatedAt                time.Time `json:"createdAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
	CancelledAt              time.Time `json:"cancelledAt"`
	ClosedAt                 time.Time `json:"closedAt"`
	Id                       string    `json:"id"`
	Name                     string    `json:"name"`
	Email                    string    `json:"email"`
	Phone                    string    `json:"phone"`
	ClientIp                 string    `json:"clientIp"`
	DisplayFinancialStatus   string    `json:"displayFinancialStatus"`
	DisplayFulfillmentStatus string    `json:"displayFulfillmentStatus"`
	CurrencyCode             string    `json:"currencyCode"`
	DiscountCodes            []string  `json:"discountCodes"`
	TotalDiscountsSet        MoneyBag  `json:"totalDiscountsSet"`
	CurrentTotalDiscountsSet MoneyBag  `json:"currentTotalDiscountsSet"`
	TotalPriceSet            MoneyBag  `json:"totalPriceSet"`
	CurrentTotalPriceSet     MoneyBag  `json:"currentTotalPriceSet"`
	SubtotalPriceSet         MoneyBag  `json:"subtotalPriceSet"`
	CurrentSubtotalPriceSet  MoneyBag  `json:"currentSubtotalPriceSet"`
	TotalRefundedSet         MoneyBag  `json:"totalRefundedSet"`
	NetPaymentSet            MoneyBag  `json:"netPaymentSet"`
	Refunds                  []Refund  `json:"refunds" graphql:"refunds(first: $first)"`
	LineItems                struct {
		Nodes []LineItem `json:"nodes"`
	} `json:"lineItems" graphql:"lineItems(first: $first)"`
	Fulfillments []Fulfillment `json:"fulfillments" graphql:"fulfillments(first: $first)"`
}

func (r *OrderRepo) List(ctx context.Context, first int, after *string) ([]Order, PageInfo, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(first),
		"after": (*graphql.String)(nil), // 默认值为 null
	}

	if after != nil {
		variables["after"] = graphql.String(*after)
	}

	var query struct {
		Orders struct {
			Nodes    []Order  `json:"nodes"`
			PageInfo PageInfo `json:"pageInfo"`
		} `json:"orders" graphql:"orders(first: $first, after: $after)"`
	}

	err := r.client.Query(ctx, &query, variables)
	return query.Orders.Nodes, query.Orders.PageInfo, err
}

func (r *OrderRepo) Get(ctx context.Context, id int64) (*Order, error) {

	variables := map[string]interface{}{
		"first": graphql.Int(10),
		"id":    fmt.Sprintf(input.FmtOrder, id),
	}

	var query struct {
		Order Order `json:"order" graphql:"order(id: $id)"`
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	if query.Order.Id == "" {
		return nil, ErrRecordNotFound
	}

	return &query.Order, nil
}

// normalizeName 归一化订单号：trim；非空且无 '#' 前缀则补 '#'（Shopify name 形如 #2735）
func normalizeName(name string) string {
	s := strings.TrimSpace(name)
	if s == "" {
		return s
	}
	if s[0] != '#' {
		s = "#" + s
	}
	return s
}

// GetByName 按 order name 搜订单（如 #2735）。
func (r *OrderRepo) GetByName(ctx context.Context, name string) (*Order, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrRecordNotFound
	}
	variables := map[string]interface{}{
		"first": graphql.Int(10), // 供 Order 内嵌 refunds/lineItems(first: $first)
		"query": graphql.String("name:" + normalizeName(name)),
	}
	var q struct {
		Orders struct {
			Nodes []Order `json:"nodes"`
		} `json:"orders" graphql:"orders(first: 1, query: $query)"`
	}
	if err := r.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}
	if len(q.Orders.Nodes) == 0 {
		return nil, ErrRecordNotFound
	}
	return &q.Orders.Nodes[0], nil
}

// AddTags 追加 tag（tagsAdd，不覆盖原有 tag）。orderGID 形如 gid://shopify/Order/123。
func (r *OrderRepo) AddTags(ctx context.Context, orderGID string, tags []string) error {
	gqlTags := make([]graphql.String, len(tags))
	for i, t := range tags {
		gqlTags[i] = graphql.String(t)
	}
	variables := map[string]interface{}{"id": orderGID, "tags": gqlTags}
	var m struct {
		TagsAdd struct {
			UserErrors UserErrors `json:"userErrors"`
		} `json:"tagsAdd" graphql:"tagsAdd(id: $id, tags: $tags)"`
	}
	if err := r.client.Mutation(ctx, &m, variables); err != nil {
		return err
	}
	for _, e := range m.TagsAdd.UserErrors {
		return errors.New(string(e.Message))
	}
	return nil
}
