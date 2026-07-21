package graphql

import (
	"context"
	"github.com/fstln-dev/utils-lib-go/pkg/goshopify/input"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
	"github.com/pkg/errors"
	"time"
)

type UsageRecordRepo struct {
	client *Client
}

func NewUsageRecordRepo(config *Client) *UsageRecordRepo {
	return &UsageRecordRepo{client: config}
}

// AppSubscription 订阅信息 用于CURD的结果解析
type AppSubscription struct {
	CreatedAt        time.Time                 `json:"createdAt"`
	CurrentPeriodEnd time.Time                 `json:"currentPeriodEnd"`
	Id               string                    `json:"id"`
	Name             string                    `json:"name"`
	ReturnUrl        string                    `json:"returnUrl"`
	Status           string                    `json:"status"`
	Test             bool                      `json:"test"`
	TrialDays        int                       `json:"trialDays"`
	LineItems        []AppSubscriptionLineItem `json:"lineItems"`
}
type AppSubscriptionLineItem struct {
	Id   string `json:"id"`
	Plan struct {
		PricingDetails struct {
			Typename            string `json:"__typename" graphql:"__typename"`
			AppRecurringPricing `graphql:"...on AppRecurringPricing"`
			AppUsagePricing     `graphql:"...on AppUsagePricing"`
		} `json:"pricingDetails"`
	} `json:"plan"`
}
type AppUsagePricing struct {
	Terms        string  `json:"terms,omitempty"`
	BalanceUsed  MoneyV2 `json:"balanceUsed,omitempty"`
	CappedAmount MoneyV2 `json:"cappedAmount,omitempty"` //上线金额
}
type AppRecurringPricing struct {
	Interval string                  `json:"interval"`
	Price    MoneyV2                 `json:"price,omitempty"`
	Discount AppSubscriptionDiscount `json:"discount,omitempty"`
}
type AppSubscriptionDiscount struct {
	//DurationLimitInIntervals 折扣可用周期数
	DurationLimitInIntervals     int     `json:"durationLimitInIntervals"`
	RemainingDurationInIntervals int     `json:"remainingDurationInIntervals"`
	PriceAfterDiscount           MoneyV2 `json:"priceAfterDiscount"`
	Value                        struct {
		Typename                      string `json:"__typename" graphql:"__typename"`
		AppSubscriptionDiscountAmount struct {
			Amount MoneyV2 `json:"amount"`
		} `graphql:"...on AppSubscriptionDiscountAmount"`
		AppSubscriptionDiscountPercentage struct {
			Percentage float64 `json:"percentage"`
		} `graphql:"...on AppSubscriptionDiscountPercentage"`
	} `json:"value"`
}

type AppUsageRecord struct {
	CreatedAt      time.Time `json:"createdAt"`
	Id             string    `json:"id"`
	IdempotencyKey string    `json:"idempotencyKey"`
	Description    string    `json:"description"`
	Price          MoneyV2   `json:"price"`
}

func (r *UsageRecordRepo) GetActiveSubscription(ctx context.Context) (*AppSubscription, error) {
	variables := map[string]interface{}{}
	var query struct {
		CurrentAppInstallation struct {
			ActiveSubscriptions []AppSubscription `json:"activeSubscriptions"`
		} `json:"currentAppInstallation" graphql:"currentAppInstallation"`
	}
	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	if len(query.CurrentAppInstallation.ActiveSubscriptions) == 0 {
		return nil, errors.Wrap(ErrRecordNotFound, "ActiveSubscriptions not exist")
	}

	return &query.CurrentAppInstallation.ActiveSubscriptions[0], nil
}

func (r *UsageRecordRepo) Create(ctx context.Context, do *input.AppUsageRecordInput) (int64, error) {
	variables := map[string]interface{}{
		"description":            graphql.String(do.Description),
		"idempotencyKey":         graphql.String(do.IdempotencyKey),
		"subscriptionLineItemId": graphql.ID(do.SubscriptionLineItemId),
		"price":                  do.Price,
	}

	var query struct {
		AppUsageRecordCreate struct {
			AppUsageRecord AppUsageRecord `json:"appUsageRecord"`
			UserErrors     UserErrors     `json:"userErrors"`
		} `json:"appUsageRecordCreate" graphql:"appUsageRecordCreate(description: $description, idempotencyKey: $idempotencyKey, subscriptionLineItemId: $subscriptionLineItemId, price: $price)"`
	}
	err := r.client.Mutation(ctx, &query, variables)
	if err != nil {
		return 0, err
	}

	for _, v := range query.AppUsageRecordCreate.UserErrors {
		return 0, errors.Wrap(ErrRecordNotFound, string(v.Message))
	}

	id := IdByGraphQL(query.AppUsageRecordCreate.AppUsageRecord.Id)
	return id, nil
}

func (r *UsageRecordRepo) Page(ctx context.Context) (*Subscription, error) {
	//TODO implement me
	panic("implement me")
}
