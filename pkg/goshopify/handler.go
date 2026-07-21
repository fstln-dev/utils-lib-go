package goshopify

import (
	"github.com/fstln-dev/utils-lib-go/pkg/goshopify/graphql"
)

// Handler manages communication with the Shopify *graphql.
type Handler struct {
	client *graphql.Client

	Shop              *graphql.ShopRepo
	Metafield         *graphql.MetafieldRepo
	Webhook           *graphql.WebhookRepo
	DiscountCodeBasic *graphql.DiscountCodeBasicRepo
	RedeemCode        *graphql.RedeemCodeRepo
	Product           *graphql.ProductRepo
	Customer          *graphql.CustomerRepo
	Segment           *graphql.SegmentRepo
	Order             *graphql.OrderRepo
	Subscription      *graphql.SubscriptionRepo
	UsageRecord       *graphql.UsageRecordRepo
}

func NewHandler(shopName, token string, opts ...graphql.Option) (*Handler, error) {
	client, err := graphql.NewClient(shopName, token, opts...)
	if err != nil {
		return nil, err
	}
	return &Handler{
		client: client,
		// graphql.Client应当复用, 在这里传入多个Repo中
		Shop:              graphql.NewShopRepo(client),
		Metafield:         graphql.NewMetafieldRepo(client),
		Webhook:           graphql.NewWebhookRepo(client),
		DiscountCodeBasic: graphql.NewDiscountCodeBasicRepo(client),
		RedeemCode:        graphql.NewRedeemCodeRepo(client),
		Product:           graphql.NewProductRepo(client),
		Customer:          graphql.NewCustomerRepo(client),
		Segment:           graphql.NewSegmentRepo(client),
		Order:             graphql.NewOrderRepo(client),
		Subscription:      graphql.NewSubscriptionRepo(client),
		UsageRecord:       graphql.NewUsageRecordRepo(client),
	}, nil
}

func WithApiVersion(version string) graphql.Option {
	return func(conf *graphql.Client) {
		conf.ApiVersion = version
	}
}
