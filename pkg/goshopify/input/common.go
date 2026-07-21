package input

const (
	FmtShop = "gid://shopify/Shop/%d"

	FmtDiscountRedeemCodeBulkCreation = "gid://shopify/DiscountRedeemCodeBulkCreation/%d"
	FmtDiscountRedeemCode             = "gid://shopify/DiscountRedeemCode/%d"
	FmtDiscountCodeNode               = "gid://shopify/DiscountCodeNode/%d"
	FmtJob                            = "gid://shopify/Job/%s"
	FmtProduct                        = "gid://shopify/Product/%d"
	FmtProductVariant                 = "gid://shopify/ProductVariant/%d"
	FmtCollection                     = "gid://shopify/Collection/%d"
	FmtCustomer                       = "gid://shopify/Customer/%d"
	FmtSegment                        = "gid://shopify/Segment/%d"
	FmtOrder                          = "gid://shopify/Order/%d"
)

type DiscountARInput struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

type MetafieldInput struct {
	Id        *int64  `json:"id,omitempty"`
	Key       *string `json:"key,omitempty"`
	Namespace *string `json:"namespace,omitempty"`
	Type      *string `json:"type,omitempty"`
	Value     *string `json:"value,omitempty"`
}

type MetafieldsSetInput struct {
	OwnerId       string  `json:"ownerId,omitempty"`
	Value         string  `json:"value,omitempty"`
	Key           string  `json:"key,omitempty"`
	Type          *string `json:"type,omitempty"`
	Namespace     *string `json:"namespace,omitempty"`
	CompareDigest *string `json:"compareDigest,omitempty"`
}

type MoneyInput struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}
