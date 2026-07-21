package graphql

import (
	"context"
	"errors"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
	"time"
)

type WebhookRepo struct {
	client *Client
}

func NewWebhookRepo(config *Client) *WebhookRepo {
	return &WebhookRepo{client: config}
}

type Webhook struct {
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Id          string    `json:"id"`
	CallbackUrl string    `json:"callbackUrl"`
	Format      string    `json:"format"`
	Topic       string    `json:"topic"`
	ApiVersion  struct {
		DisplayName string `json:"displayName"`
		Handle      string `json:"handle"`
		Supported   bool   `json:"supported"`
	} `json:"apiVersion"`
}

func (r *WebhookRepo) Create(ctx context.Context, topic, callbackUrl string) (*Webhook, error) {
	type WebhookSubscriptionTopic string
	type WebhookSubscriptionInput struct {
		CallbackUrl string `json:"callbackUrl"`
	}
	variables := map[string]interface{}{
		"topic": WebhookSubscriptionTopic(topic),
		"webhookSubscription": WebhookSubscriptionInput{
			CallbackUrl: callbackUrl,
		},
	}
	var query struct {
		WebhookSubscriptions struct {
			WebhookSubscription Webhook    `json:"webhookSubscription"`
			UserErrors          UserErrors `json:"userErrors"`
		} `json:"webhookSubscriptionCreate" graphql:"webhookSubscriptionCreate(topic: $topic, webhookSubscription: $webhookSubscription)"`
	}

	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.WebhookSubscriptions.UserErrors {
		//todo wrap error
		return nil, errors.New(string(v.Message))
	}

	return &query.WebhookSubscriptions.WebhookSubscription, nil
}

func (r *WebhookRepo) List(ctx context.Context, first int) ([]Webhook, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(first),
	}

	var query struct {
		WebhookSubscriptionsObj struct {
			Nodes    []Webhook `json:"nodes"`
			PageInfo PageInfo  `json:"pageInfo"`
		} `json:"webhookSubscriptions" graphql:"webhookSubscriptions(first: $first)"`
	}

	err := r.client.Query(ctx, &query, variables)
	return query.WebhookSubscriptionsObj.Nodes, err
}
