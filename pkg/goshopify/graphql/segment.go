package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/fstln-dev/utils-lib-go/pkg/goshopify/input"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
	"time"
)

type SegmentRepo struct {
	client *Client
}

func NewSegmentRepo(config *Client) *SegmentRepo {
	return &SegmentRepo{client: config}
}

type CustomerSegment struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Query        string    `json:"query"`
	CreationDate time.Time `json:"creationDate"`
	LastEditDate time.Time `json:"lastEditDate"`
}

func (r *SegmentRepo) Create(ctx context.Context, name, queryStr string) (*CustomerSegment, error) {
	variables := map[string]interface{}{
		"name":  graphql.String(name),
		"query": graphql.String(queryStr),
	}
	var query struct {
		SegmentCreate struct {
			Segment    CustomerSegment `json:"Segment"`
			UserErrors UserErrors      `json:"userErrors"`
		} `json:"segmentCreate" graphql:"segmentCreate(name: $name, query: $query)"`
	}

	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.SegmentCreate.UserErrors {
		//todo wrap error
		return nil, errors.New(string(v.Message))
	}

	return &query.SegmentCreate.Segment, nil
}

func (r *SegmentRepo) List(ctx context.Context, first int) ([]CustomerSegment, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(first),
	}

	var query struct {
		Segments struct {
			Nodes    []CustomerSegment `json:"nodes"`
			PageInfo PageInfo          `json:"pageInfo"`
		} `json:"segments" graphql:"segments(first: $first)"`
	}

	err := r.client.Query(ctx, &query, variables)
	return query.Segments.Nodes, err
}

func (r *SegmentRepo) GetByID(ctx context.Context, id int64) (*CustomerSegment, error) {
	variables := graphql.Variables{
		"id": fmt.Sprintf(input.FmtSegment, id),
	}
	var query struct {
		CustomerSegment `json:"segment" graphql:"segment(id: $id)"`
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	if query.CustomerSegment.Id == "" {
		return nil, ErrRecordNotFound
	}

	return &query.CustomerSegment, nil
}

func (r *SegmentRepo) Update(ctx context.Context, id int64, name, queryStr string) (*CustomerSegment, error) {
	variables := graphql.Variables{
		"id":    fmt.Sprintf(input.FmtSegment, id),
		"name":  graphql.String(name),
		"query": graphql.String(queryStr),
	}

	var query struct {
		SegmentUpdate struct {
			Segment    CustomerSegment `json:"segment"`
			UserErrors UserErrors      `json:"userErrors"`
		} `json:"segmentUpdate" graphql:"segmentUpdate(id: $id, name: $name, query: $query)"`
	}

	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.SegmentUpdate.UserErrors {
		//todo wrap error
		return nil, errors.New(string(v.Message))
	}

	return &query.SegmentUpdate.Segment, nil
}

func (r *SegmentRepo) Delete(ctx context.Context, id int64) error {
	variables := graphql.Variables{
		"id": fmt.Sprintf(input.FmtSegment, id),
	}
	var query struct {
		SegmentDelete struct {
			DeletedSegmentId graphql.ID `json:"deletedSegmentId"`
			UserErrors       UserErrors `json:"userErrors"`
		} `json:"segmentDelete" graphql:"segmentDelete(id: $id)"`
	}

	return r.client.Mutation(ctx, &query, variables)
}
