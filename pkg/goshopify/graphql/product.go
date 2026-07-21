package graphql

import (
	"context"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
	"time"
)

type ProductRepo struct {
	client *Client
}

func NewProductRepo(config *Client) *ProductRepo {
	return &ProductRepo{client: config}
}

type Variant struct {
	Id               string `json:"id"`
	Title            string `json:"title"`
	Price            string `json:"price"`
	Sku              string `json:"sku"`
	AvailableForSale bool   `json:"availableForSale"`
}
type Image struct {
	Id          string `json:"id"`
	OriginalSrc string `json:"originalSrc"`
	AltText     string `json:"altText"`
}
type Product struct {
	Id              string    `json:"id"`
	Title           string    `json:"title"`
	DescriptionHtml string    `json:"descriptionHtml"`
	Vendor          string    `json:"vendor"`
	ProductType     string    `json:"productType"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	VariantsObj     struct {
		Nodes    []Variant `json:"nodes"`
		PageInfo PageInfo  `json:"pageInfo"`
	} `json:"variants" graphql:"variants(first: $first)"`
	ImagesObj struct {
		Nodes    []Image  `json:"nodes"`
		PageInfo PageInfo `json:"pageInfo"`
	} `json:"images" graphql:"images(first: $first)"`
}
type Collection struct {
	UpdatedAt       time.Time `json:"updatedAt"`
	Id              string    `json:"id"`
	Handle          string    `json:"handle"`
	Title           string    `json:"title"`
	SortOrder       string    `json:"sortOrder"`
	TemplateSuffix  string    `json:"templateSuffix"`
	DescriptionHtml string    `json:"descriptionHtml"`
}

func (r *ProductRepo) List(ctx context.Context, first int) ([]Product, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(first),
	}

	var query struct {
		Products struct {
			Nodes    []Product `json:"nodes"`
			PageInfo PageInfo  `json:"pageInfo"`
		} `json:"products" graphql:"products(first: $first)"`
	}

	err := r.client.Query(ctx, &query, variables)
	return query.Products.Nodes, err
}

func (r *ProductRepo) ListCollection(ctx context.Context, first int) ([]Collection, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(first),
	}

	var query struct {
		Collections struct {
			Nodes    []Collection `json:"nodes"`
			PageInfo PageInfo     `json:"pageInfo"`
		} `json:"collections" graphql:"collections(first: $first)"`
	}

	err := r.client.Query(ctx, &query, variables)
	return query.Collections.Nodes, err
}
