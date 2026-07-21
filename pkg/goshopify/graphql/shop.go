package graphql

import (
	"context"
	"time"
)

type ShopRepo struct {
	client *Client
}

func NewShopRepo(client *Client) *ShopRepo {
	return &ShopRepo{client: client}
}

type Shop struct {
	CreatedAt                time.Time `json:"createdAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
	Id                       string    `json:"id"`
	Name                     string    `json:"name"`
	MyshopifyDomain          string    `json:"myshopifyDomain"`
	Url                      string    `json:"url"`
	Email                    string    `json:"email"`
	ContactEmail             string    `json:"contactEmail"`
	CurrencyCode             string    `json:"currencyCode"`
	CustomerAccounts         string    `json:"customerAccounts"`
	Description              string    `json:"description"`
	StorefrontUrl            string    `json:"storefrontUrl"`
	IanaTimezone             string    `json:"ianaTimezone"`
	TimezoneAbbreviation     string    `json:"timezoneAbbreviation"`
	TimezoneOffset           string    `json:"timezoneOffset"`
	TimezoneOffsetMinutes    int       `json:"timezoneOffsetMinutes"`
	PublicationCount         int       `json:"publicationCount"`
	TransactionalSmsDisabled bool      `json:"transactionalSmsDisabled"`
	TaxShipping              bool      `json:"taxShipping"`
	TaxesIncluded            bool      `json:"taxesIncluded"`
	CurrencyFormats          struct {
		MoneyFormat                     string `json:"moneyFormat"`
		MoneyInEmailsFormat             string `json:"moneyInEmailsFormat"`
		MoneyWithCurrencyFormat         string `json:"moneyWithCurrencyFormat"`
		MoneyWithCurrencyInEmailsFormat string `json:"moneyWithCurrencyInEmailsFormat"`
	} `json:"currencyFormats"`
}

func (r *ShopRepo) Get(ctx context.Context) (*Shop, error) {
	var query struct {
		Shop Shop `json:"shop" graphql:"shop"`
	}

	err := r.client.Query(ctx, &query, nil)
	return &query.Shop, err
}
