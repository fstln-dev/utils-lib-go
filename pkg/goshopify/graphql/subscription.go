package graphql

import (
	"context"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
)

type SubscriptionRepo struct {
	client *Client
}

func NewSubscriptionRepo(config *Client) *SubscriptionRepo {
	return &SubscriptionRepo{client: config}
}

type Subscription struct {
}

func (r *SubscriptionRepo) GetActive(ctx context.Context) (*Subscription, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(100000),
	}

	var query struct {
		SubscriptionsObj struct {
			Nodes    []Subscription `json:"nodes"`
			PageInfo PageInfo       `json:"pageInfo"`
		} `json:"subscription" graphql:"subscription"`
	}

	if err := r.client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *SubscriptionRepo) CancelActive(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
