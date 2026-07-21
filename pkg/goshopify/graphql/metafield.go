package graphql

import (
	"context"
	"errors"
	"github.com/fstln-dev/utils-lib-go/pkg/goshopify/input"
	"time"
)

type MetafieldRepo struct {
	client *Client
}

func NewMetafieldRepo(config *Client) *MetafieldRepo {
	return &MetafieldRepo{client: config}
}

type MetafieldDefinition struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Key         string `json:"key,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	Description string `json:"description,omitempty"`
}

type Metafield struct {
	CreatedAt     time.Time            `json:"createdAt,omitempty"`
	UpdatedAt     time.Time            `json:"updatedAt,omitempty"`
	ID            string               `json:"id,omitempty"`
	Key           string               `json:"key,omitempty"`
	Value         string               `json:"value,omitempty"`
	Type          string               `json:"type,omitempty"`
	Namespace     string               `json:"namespace,omitempty"`
	OwnerType     string               `json:"ownerType,omitempty"`
	Description   string               `json:"description,omitempty"`
	CompareDigest string               `json:"compareDigest,omitempty"`
	Definition    *MetafieldDefinition `json:"definition,omitempty"`
	//JsonValue   interface{} `json:"jsonValue,omitempty"` //todo 待测试
}

func (r *MetafieldRepo) Set(ctx context.Context, ins []input.MetafieldsSetInput) ([]Metafield, error) {
	variables := map[string]interface{}{"metafields": ins}
	var query struct {
		MetafieldsSet struct {
			Metafields []Metafield `json:"metafields"`
			UserErrors UserErrors  `json:"userErrors"`
		} `json:"metafieldSubscriptionCreate" graphql:"metafieldsSet(metafields: $metafields)"`
	}

	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.MetafieldsSet.UserErrors {
		//todo wrap error
		return nil, errors.New(string(v.Message))
	}

	return query.MetafieldsSet.Metafields, nil
}
