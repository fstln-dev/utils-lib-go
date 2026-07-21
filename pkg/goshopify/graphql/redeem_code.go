package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/fstln-dev/utils-lib-go/pkg/goshopify/input"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
	"time"
)

type RedeemCodeRepo struct {
	client *Client
}

func NewRedeemCodeRepo(client *Client) *RedeemCodeRepo {
	return &RedeemCodeRepo{client: client}
}

type DiscountRedeemCodeBulkCreation struct {
	CreatedAt     time.Time `json:"createdAt"`
	Id            string    `json:"id"`
	Done          bool      `json:"done"`
	CodesCount    int       `json:"codesCount"`
	FailedCount   int       `json:"failedCount"`
	ImportedCount int       `json:"importedCount"`
	Codes         struct {
		Nodes    []DiscountRedeemCodeBulkCreationCode `json:"nodes"`
		PageInfo PageInfo                             `json:"pageInfo"`
	} `json:"codes" graphql:"codes(first: $first)"`
}

type DiscountRedeemCodeBulkCreationCode struct {
	Code               string              `json:"code"`
	DiscountRedeemCode *DiscountRedeemCode `json:"discountRedeemCode"`
	Errors             []DiscountUserError `json:"errors"`
}

type DiscountUserError struct {
	Code      string
	Message   string
	ExtraInfo string
	Field     []string
}

// CodesConnection 代表折扣码的分页连接，包含边缘和分页信息。
type CodesConnection struct {
	PageInfo PageInfo           `json:"pageInfo"`
	Edges    []DiscountCodeEdge `json:"edges"`
}

// DiscountCodeEdge 代表分页查询中的一个边缘，包含游标和节点。
type DiscountCodeEdge struct {
	Cursor string             `json:"cursor"`
	Node   DiscountRedeemCode `json:"node"`
}

type PaginatedDiscountCodes struct {
	Codes    []DiscountRedeemCode
	PageInfo PageInfo
}

type DiscountRedeemCode struct {
	AsyncUsageCount int    `json:"asyncUsageCount"`
	Code            string `json:"code"`
	Id              string `json:"id"`
}

func (r *RedeemCodeRepo) CreateCreation(
	ctx context.Context,
	discountId int64,
	codes []string,
) (*DiscountRedeemCodeBulkCreation, error) {
	type DiscountRedeemCodeInput struct {
		Code string `json:"code"`
	}
	cis := make([]DiscountRedeemCodeInput, len(codes))
	for i, code := range codes {
		cis[i] = DiscountRedeemCodeInput{Code: code}
	}
	variables := map[string]interface{}{
		"discountId": fmt.Sprintf(input.FmtDiscountCodeNode, discountId),
		"codes":      cis,
		"first":      graphql.Int(100),
	}

	var query struct {
		DiscountRedeemCodeBulkAdd struct {
			BulkCreation DiscountRedeemCodeBulkCreation `json:"bulkCreation"`
			UserErrors   UserErrors                     `json:"userErrors"`
		} `json:"discountRedeemCodeBulkAdd" graphql:"discountRedeemCodeBulkAdd(discountId: $discountId, codes: $codes)"`
	}
	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.DiscountRedeemCodeBulkAdd.UserErrors {
		return nil, errors.New(string(v.Message))
	}

	return &query.DiscountRedeemCodeBulkAdd.BulkCreation, nil
}

func (r *RedeemCodeRepo) GetCreation(ctx context.Context, id int64) (*DiscountRedeemCodeBulkCreation, error) {
	variables := map[string]interface{}{
		"id":    fmt.Sprintf(input.FmtDiscountRedeemCodeBulkCreation, id),
		"first": graphql.Int(100),
	}

	var query struct {
		DiscountRedeemCodeBulkCreation `json:"discountRedeemCodeBulkCreation" graphql:"discountRedeemCodeBulkCreation(id: $id)"`
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	if query.DiscountRedeemCodeBulkCreation.Id == "" {
		return nil, ErrRecordNotFound
	}

	return &query.DiscountRedeemCodeBulkCreation, err
}

func (r *RedeemCodeRepo) BulkDel(ctx context.Context, discountId int64, redeemGids []string) (string, error) {
	variables := map[string]interface{}{
		"discountId": fmt.Sprintf(input.FmtDiscountCodeNode, discountId),
		"ids":        redeemGids,
	}

	var query struct {
		DiscountCodeRedeemCodeBulkDelete struct {
			Job        Job        `json:"job"`
			UserErrors UserErrors `json:"userErrors"`
		} `json:"discountCodeRedeemCodeBulkDelete" graphql:"discountCodeRedeemCodeBulkDelete(discountId: $discountId, ids: $ids)"`
	}
	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return "", err
	}

	for _, v := range query.DiscountCodeRedeemCodeBulkDelete.UserErrors {
		return "", errors.New(string(v.Message))
	}

	return query.DiscountCodeRedeemCodeBulkDelete.Job.Id, nil
}

func (r *RedeemCodeRepo) ListDiscountCodesWithPagination(ctx context.Context, discountId int64, pageSize int, after *string) (*PaginatedDiscountCodes, error) {
	if pageSize > 250 || pageSize <= 0 {
		pageSize = 250
	}

	// 准备 GraphQL 查询的变量。
	variables := map[string]interface{}{
		"discountNodeId": fmt.Sprintf(input.FmtDiscountCodeNode, discountId),
		"first":          graphql.Int(pageSize),
		"after":          (*graphql.String)(nil), // 默认值为 null
	}

	if after != nil {
		variables["after"] = graphql.String(*after)
	}

	var query struct {
		Node struct {
			DiscountCodeNode struct {
				CodeDiscount struct {
					DiscountCodeBasic struct {
						Codes CodesConnection `graphql:"codes(first: $first, after: $after)"`
					} `graphql:"... on DiscountCodeBasic"`
				} `json:"codeDiscount"`
			} `graphql:"... on DiscountCodeNode"`
		} `graphql:"node(id: $discountNodeId)"`
	}

	// 执行查询
	if err := r.client.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("GraphQL query failed: %w", err)
	}

	// 从嵌套的查询结果中提取数据
	resultConnection := query.Node.DiscountCodeNode.CodeDiscount.DiscountCodeBasic.Codes

	// 将 GraphQL 的 'edges' 结构扁平化为简单的折扣码切片，方便调用者使用。
	codes := make([]DiscountRedeemCode, 0, len(resultConnection.Edges))
	for _, edge := range resultConnection.Edges {
		codes = append(codes, edge.Node)
	}

	// 构造并返回最终结果
	paginatedResult := &PaginatedDiscountCodes{
		Codes:    codes,
		PageInfo: resultConnection.PageInfo,
	}

	return paginatedResult, nil
}

func (r *RedeemCodeRepo) GetJobByID(ctx context.Context, jobId string) (*Job, error) {
	variables := map[string]interface{}{
		"id": fmt.Sprintf(input.FmtJob, jobId),
	}

	var query struct {
		Job Job `json:"job" graphql:"job(id: $id)"`
	}

	if err := r.client.Query(ctx, &query, variables); err != nil {
		return nil, err
	}

	return &query.Job, nil
}

// ListRedeem todo
func (r *RedeemCodeRepo) ListRedeem(ctx context.Context, id int64) (*DiscountRedeemCodeBulkCreation, error) {
	return nil, nil
}
