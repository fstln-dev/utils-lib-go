package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/fstln-dev/utils-lib-go/pkg/goshopify/input"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
	"time"
)

type DiscountCodeBasicRepo struct {
	client *Client
}

func NewDiscountCodeBasicRepo(client *Client) *DiscountCodeBasicRepo {
	return &DiscountCodeBasicRepo{client: client}
}

type DiscountCodeBasic struct {
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	StartsAt               time.Time `json:"startsAt"`
	EndsAt                 time.Time `json:"endsAt"`
	Title                  string    `json:"title"`
	Summary                string    `json:"summary"`
	ShortSummary           string    `json:"shortSummary"`
	Status                 string    `json:"status"`
	DiscountClass          string    `json:"discountClass"`
	AsyncUsageCount        int       `json:"asyncUsageCount"`
	UsageLimit             int       `json:"usageLimit"`
	RecurringCycleLimit    int       `json:"recurringCycleLimit"`
	AppliesOncePerCustomer bool      `json:"appliesOncePerCustomer"`
	HasTimelineComment     bool      `json:"hasTimelineComment"`
	CodesCount             struct {
		Count     int    `json:"count"`
		Precision string `json:"precision"`
	} `json:"codesCount"`
	CombinesWith struct {
		OrderDiscounts    bool `json:"orderDiscounts"`
		ProductDiscounts  bool `json:"productDiscounts"`
		ShippingDiscounts bool `json:"shippingDiscounts"`
	} `json:"combinesWith"`
	ShareableUrls []struct {
		Title string `json:"title"`
		Url   string `json:"url"`
	} `json:"shareableUrls"`
	TotalSales         MoneyV2                    `json:"totalSales"`
	MinimumRequirement DiscountMinimumRequirement `json:"minimumRequirement"`
	CustomerSelection  DiscountCustomerSelection  `json:"customerSelection"`
	CustomerGets       DiscountCustomerGets       `json:"customerGets"`
}

type DiscountMinimumRequirement struct {
	DiscountMinimumQuantity struct {
		Typename                     string `json:"__typename" graphql:"__typename"`
		GreaterThanOrEqualToQuantity string `json:"greaterThanOrEqualToQuantity"`
	} `graphql:"... on DiscountMinimumQuantity"`
	DiscountMinimumSubtotal struct {
		Typename                     string  `json:"__typename" graphql:"__typename"`
		GreaterThanOrEqualToSubtotal MoneyV2 `json:"greaterThanOrEqualToSubtotal"`
	} `graphql:"... on DiscountMinimumSubtotal"`
}

type DiscountCustomerSelection struct {
	DiscountCustomerAll struct {
		Typename     string `json:"__typename" graphql:"__typename"`
		AllCustomers bool   `json:"allCustomers"`
	} `graphql:"... on DiscountCustomerAll"`
	DiscountCustomerSegments struct {
		Typename string            `json:"__typename" graphql:"__typename"`
		Segments []CustomerSegment `json:"segments"`
	} `graphql:"... on DiscountCustomerSegments"`
	DiscountCustomers struct {
		Typename  string     `json:"__typename" graphql:"__typename"`
		Customers []Customer `json:"customers"`
	} `graphql:"... on DiscountCustomers"`
}

type DiscountCustomerGets struct {
	AppliesOnOneTimePurchase bool                      `json:"appliesOnOneTimePurchase"`
	AppliesOnSubscription    bool                      `json:"appliesOnSubscription"`
	Items                    DiscountItems             `json:"items"`
	Value                    DiscountCustomerGetsValue `json:"value"`
}

type DiscountItems struct {
	AllDiscountItems struct {
		Typename string `json:"__typename" graphql:"__typename"`
		AllItems bool   `json:"allItems"`
	} `graphql:"... on AllDiscountItems"`
	DiscountCollections struct {
		Typename    string `json:"__typename" graphql:"__typename"`
		Collections struct {
			Nodes []Collection `json:"collections"`
		} `json:"collections" graphql:"collections(first: $first)"`
	} `graphql:"... on DiscountCollections"`
	DiscountProducts struct {
		Typename string `json:"__typename" graphql:"__typename"`
		Products struct {
			Nodes []Product `json:"nodes"`
		} `json:"products" graphql:"products(first: $first)"`
		ProductVariants struct {
			Nodes []Variant `json:"nodes"`
		} `json:"productVariants" graphql:"productVariants(first: $first)"`
	} `graphql:"... on DiscountProducts"`
}

type DiscountCustomerGetsValue struct {
	DiscountAmount     `graphql:"... on DiscountAmount"`
	DiscountPercentage `graphql:"... on DiscountPercentage"`
	DiscountOnQuantity struct {
		Typename string `json:"__typename" graphql:"__typename"`
		Quantity struct {
			Quantity string `json:"quantity"`
		} `json:"quantity"`
		Effect struct {
			Typename           string `json:"__typename" graphql:"__typename"`
			DiscountAmount     `graphql:"... on DiscountAmount"`
			DiscountPercentage `graphql:"... on DiscountPercentage"`
		}
	} `graphql:"... on DiscountOnQuantity"`
}

type DiscountAmount struct {
	Typename          string  `json:"__typename" graphql:"__typename"`
	AppliesOnEachItem bool    `json:"appliesOnEachItem"`
	Amount            MoneyV2 `json:"amount"`
}

type DiscountPercentage struct {
	Typename   string  `json:"__typename" graphql:"__typename"`
	Percentage float64 `json:"percentage"`
}

type DiscountCodeNode struct {
	Id           string `json:"id"`
	CodeDiscount struct {
		Typename          string `json:"__typename" graphql:"__typename"`
		DiscountCodeBasic `graphql:"... on DiscountCodeBasic"`
	} `json:"codeDiscount"`
}

func (r *DiscountCodeBasicRepo) Create(
	ctx context.Context,
	basicCodeDiscount *input.DiscountCodeBasicInput,
) (*DiscountCodeNode, error) {
	variables := map[string]interface{}{
		"first":             graphql.Int(10),
		"basicCodeDiscount": *basicCodeDiscount,
	}

	var query struct {
		DiscountCodeBasicCreate struct {
			CodeDiscountNode DiscountCodeNode `json:"codeDiscountNode"`
			UserErrors       UserErrors       `json:"userErrors"`
		} `json:"discountCodeBasicCreate" graphql:"discountCodeBasicCreate(basicCodeDiscount: $basicCodeDiscount)"`
	}
	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.DiscountCodeBasicCreate.UserErrors {
		return nil, errors.New(string(v.Message))
	}

	return &query.DiscountCodeBasicCreate.CodeDiscountNode, nil
}

func (r *DiscountCodeBasicRepo) Update(
	ctx context.Context,
	id int64,
	basicCodeDiscount *input.DiscountCodeBasicInput,
) (*DiscountCodeNode, error) {
	variables := map[string]interface{}{
		"first":             graphql.Int(10),
		"id":                fmt.Sprintf(input.FmtDiscountCodeNode, id),
		"basicCodeDiscount": *basicCodeDiscount,
	}
	var query struct {
		DiscountCodeBasicUpdate struct {
			CodeDiscountNode DiscountCodeNode `json:"codeDiscountNode"`
			UserErrors       UserErrors       `json:"userErrors"`
		} `json:"discountCodeBasicUpdate" graphql:"discountCodeBasicUpdate(id: $id, basicCodeDiscount: $basicCodeDiscount)"`
	}
	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.DiscountCodeBasicUpdate.UserErrors {
		return nil, errors.New(string(v.Message))
	}

	return &query.DiscountCodeBasicUpdate.CodeDiscountNode, nil
}

type DiscountNode struct {
	Id       string `json:"id"`
	Discount struct {
		Typename          string `json:"__typename" graphql:"__typename"`
		DiscountCodeBasic `graphql:"... on DiscountCodeBasic"`
	} `json:"discount"`
}

func (r *DiscountCodeBasicRepo) List(ctx context.Context, first int) ([]DiscountNode, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(first),
	}

	var query struct {
		DiscountNodes struct {
			Nodes    []DiscountNode `json:"nodes"`
			PageInfo PageInfo       `json:"pageInfo"`
		} `json:"discountNodes" graphql:"discountNodes(first: $first)"`
	}

	err := r.client.Query(ctx, &query, variables)
	return query.DiscountNodes.Nodes, err
}

func (r *DiscountCodeBasicRepo) Get(ctx context.Context, id int64) (*DiscountNode, error) {
	variables := map[string]interface{}{
		"id":    fmt.Sprintf(input.FmtDiscountCodeNode, id),
		"first": graphql.Int(10),
	}

	var query struct {
		DiscountNode `json:"discountNode" graphql:"discountNode(id: $id)"`
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	if query.DiscountNode.Id == "" {
		return nil, ErrRecordNotFound
	}

	return &query.DiscountNode, err
}
