package input

type AppUsageRecordInput struct {
	SubscriptionLineItemId string     `json:"subscriptionLineItemId"` //订阅行ID
	IdempotencyKey         string     `json:"idempotencyKey"`         //客户唯一标识, 避免重复收费
	Description            string     `json:"description"`
	Price                  MoneyInput `json:"price"`
}
