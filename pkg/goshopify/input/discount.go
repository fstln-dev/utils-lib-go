package input

import "time"

type DiscountCodeBasicInput struct {
	Code                   *string                          `json:"code,omitempty"`
	Title                  *string                          `json:"title,omitempty"`
	AppliesOncePerCustomer *bool                            `json:"appliesOncePerCustomer,omitempty"`
	RecurringCycleLimit    *int                             `json:"recurringCycleLimit,omitempty"`
	UsageLimit             *int                             `json:"usageLimit,omitempty"`
	StartsAt               *time.Time                       `json:"startsAt,omitempty"`
	EndsAt                 *time.Time                       `json:"endsAt"`
	CombinesWith           *DiscountCombinesWithInput       `json:"combinesWith,omitempty"`
	MinimumRequirement     *DiscountMinimumRequirementInput `json:"minimumRequirement,omitempty"`
	CustomerSelection      *DiscountCustomerSelectionInput  `json:"customerSelection,omitempty"`
	CustomerGets           *DiscountCustomerGetsInput       `json:"customerGets,omitempty"`
}

type DiscountCombinesWithInput struct {
	OrderDiscounts    *bool `json:"orderDiscounts,omitempty"`
	ProductDiscounts  *bool `json:"productDiscounts,omitempty"`
	ShippingDiscounts *bool `json:"shippingDiscounts,omitempty"`
}

type DiscountMinimumRequirementInput struct {
	Quantity *DiscountMinimumQuantityInput `json:"quantity"`
	Subtotal *DiscountMinimumSubtotalInput `json:"subtotal"`
}

type DiscountMinimumQuantityInput struct {
	GreaterThanOrEqualToQuantity *string `json:"greaterThanOrEqualToQuantity"`
}

type DiscountMinimumSubtotalInput struct {
	GreaterThanOrEqualToSubtotal *string `json:"greaterThanOrEqualToSubtotal"`
}

type DiscountCustomerSelectionInput struct {
	AllCustomers     *bool            `json:"all"`
	CustomerSegments *DiscountARInput `json:"customerSegments"`
	Customers        *DiscountARInput `json:"customers"`
}

type DiscountCustomerGetsInput struct {
	AppliesOnOneTimePurchase *bool                           `json:"appliesOnOneTimePurchase,omitempty"`
	AppliesOnSubscription    *bool                           `json:"appliesOnSubscription,omitempty"`
	Items                    *DiscountItemsInput             `json:"items,omitempty"`
	Value                    *DiscountCustomerGetsValueInput `json:"value,omitempty"`
}

type DiscountItemsInput struct {
	All         *bool                  `json:"all"`
	Collections *DiscountARInput       `json:"collections"`
	Products    *DiscountProductsInput `json:"products"`
}

type DiscountProductsInput struct {
	ProductsToAdd           []string `json:"productsToAdd,omitempty"`
	ProductVariantsToAdd    []string `json:"productVariantsToAdd,omitempty"`
	ProductsToRemove        []string `json:"productsToRemove,omitempty"`
	ProductVariantsToRemove []string `json:"productVariantsToRemove,omitempty"`
}

type DiscountCustomerGetsValueInput struct {
	Percentage         *float64                 `json:"percentage,omitempty"`
	DiscountAmount     *DiscountAmountInput     `json:"discountAmount,omitempty"`
	DiscountOnQuantity *DiscountOnQuantityInput `json:"discountOnQuantity,omitempty"`
}

type DiscountAmountInput struct {
	Amount            *string `json:"amount,omitempty"`
	AppliesOnEachItem *bool   `json:"appliesOnEachItem,omitempty"`
}

type DiscountOnQuantityInput struct {
	Quantity *string              `json:"quantity,omitempty"`
	Effect   *DiscountEffectInput `json:"effect,omitempty"`
}

type DiscountEffectInput struct {
	Amount     *string  `json:"Amount,omitempty"`
	Percentage *float64 `json:"percentage,omitempty"`
}
