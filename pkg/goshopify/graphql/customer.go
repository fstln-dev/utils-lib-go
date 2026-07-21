package graphql

import (
	"context"
	"errors"
	"fmt"
	"github.com/fstln-dev/utils-lib-go/pkg/goshopify/input"
	"github.com/fstln-dev/utils-lib-go/pkg/graphql"
	"time"
)

type CustomerRepo struct {
	client *Client
}

func NewCustomerRepo(config *Client) *CustomerRepo {
	return &CustomerRepo{client: config}
}

type MailingAddress struct {
	Address1                string   `json:"address1,omitempty"`
	Address2                string   `json:"address2,omitempty"`
	City                    string   `json:"city,omitempty"`
	Company                 string   `json:"company,omitempty"`
	Country                 string   `json:"country,omitempty"`
	CountryCodeV2           string   `json:"countryCode,omitempty"`
	FirstName               string   `json:"firstName,omitempty"`
	LastName                string   `json:"lastName,omitempty"`
	Phone                   string   `json:"phone,omitempty"`
	ProvinceCode            string   `json:"provinceCode,omitempty"`
	Zip                     string   `json:"zip,omitempty"`
	ID                      string   `json:"id,omitempty"`
	Name                    string   `json:"name,omitempty"`
	Province                string   `json:"province,omitempty"`
	TimeZone                string   `json:"timeZone,omitempty"`
	ValidationResultSummary string   `json:"validationResultSummary,omitempty"`
	Latitude                float64  `json:"latitude,omitempty"`
	Longitude               float64  `json:"longitude,omitempty"`
	CoordinatesValidated    bool     `json:"coordinatesValidated,omitempty"`
	FormattedArea           string   `json:"formattedArea,omitempty"`
	Formatted               []string `json:"formatted,omitempty"`
}

type Customer struct {
	CanDelete               bool      `json:"canDelete"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
	LegacyResourceId        string    `json:"legacyResourceId"`
	Id                      string    `json:"id"`
	DisplayName             string    `json:"displayName"`
	FirstName               string    `json:"firstName"`
	LastName                string    `json:"lastName"`
	Email                   string    `json:"email"`
	ValidEmailAddress       bool      `json:"validEmailAddress"`
	VerifiedEmail           bool      `json:"verifiedEmail"`
	Phone                   string    `json:"phone"`
	LifetimeDuration        string    `json:"lifetimeDuration"`
	Locale                  string    `json:"locale"`
	Note                    string    `json:"note"`
	NumberOfOrders          string    `json:"numberOfOrders"`
	State                   string    `json:"state"`
	UnsubscribeUrl          string    `json:"unsubscribeUrl"`
	ProductSubscriberStatus string    `json:"productSubscriberStatus"`
	Tags                    []string  `json:"tags"`
	SmsMarketingConsent     struct {
		ConsentCollectedFrom string    `json:"consentCollectedFrom"`
		ConsentUpdatedAt     time.Time `json:"consentUpdatedAt"`
		MarketingOptInLevel  string    `json:"marketingOptInLevel"`
		MarketingState       string    `json:"marketingState"`
	} `json:"smsMarketingConsent"`
	EmailMarketingConsent struct {
		ConsentUpdatedAt    time.Time `json:"consentUpdatedAt"`
		MarketingOptInLevel string    `json:"marketingOptInLevel"`
		MarketingState      string    `json:"marketingState"`
	} `json:"emailMarketingConsent"`
	AddressesV2 struct {
		Nodes []MailingAddress `json:"nodes"`
	} `json:"addressesV2" graphql:"addressesV2(first: $first)"`
}

type CustomerSimple struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
}

func (r *CustomerRepo) List(ctx context.Context, first int, queryStr string) ([]Customer, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(first),
		"query": graphql.String(queryStr),
	}

	var query struct {
		Customers struct {
			Nodes    []Customer `json:"nodes"`
			PageInfo PageInfo   `json:"pageInfo"`
		} `json:"customers" graphql:"customers(first: $first, query: $query)"`
	}

	err := r.client.Query(ctx, &query, variables)
	return query.Customers.Nodes, err
}

func (r *CustomerRepo) GetSimple(ctx context.Context, id int64) (*CustomerSimple, error) {
	variables := map[string]interface{}{
		"id": fmt.Sprintf(input.FmtCustomer, id),
	}

	var query struct {
		Customer CustomerSimple `json:"customer" graphql:"customer(id: $id)"`
	}

	err := r.client.Query(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	if query.Customer.Id == "" {
		return nil, ErrRecordNotFound
	}

	return &query.Customer, err
}

func (r *CustomerRepo) Create(ctx context.Context, in *input.CustomerInput) (*Customer, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(10),
		"input": *in,
	}

	var query struct {
		CustomerCreate struct {
			Customer   Customer   `json:"customer"`
			UserErrors UserErrors `json:"userErrors"`
		} `json:"customerCreate" graphql:"customerCreate(input: $input)"`
	}
	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.CustomerCreate.UserErrors {
		return nil, errors.New(string(v.Message))
	}

	return &query.CustomerCreate.Customer, nil
}

func (r *CustomerRepo) Update(ctx context.Context, in *input.CustomerInput) (*Customer, error) {
	variables := map[string]interface{}{
		"first": graphql.Int(10),
		"input": *in,
	}
	var query struct {
		CustomerUpdate struct {
			Customer   Customer   `json:"customer"`
			UserErrors UserErrors `json:"userErrors"`
		} `json:"customerUpdate" graphql:"customerUpdate(input: $input)"`
	}
	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.CustomerUpdate.UserErrors {
		return nil, errors.New(string(v.Message))
	}

	return &query.CustomerUpdate.Customer, nil
}

func (r *CustomerRepo) UpdateEmailMarketing(ctx context.Context, id int64, in *input.CustomerMarketingConsentInput) (*Customer, error) {
	type CustomerEmailMarketingConsentUpdateInput struct {
		CustomerId            string                              `json:"customerId"`
		EmailMarketingConsent input.CustomerMarketingConsentInput `json:"emailMarketingConsent,omitempty"`
	}

	variables := map[string]interface{}{
		"first": graphql.Int(10),
		"input": CustomerEmailMarketingConsentUpdateInput{
			CustomerId:            fmt.Sprintf(input.FmtCustomer, id),
			EmailMarketingConsent: *in,
		},
	}
	var query struct {
		UpdateEmailMarketing struct {
			Customer   Customer   `json:"customer"`
			UserErrors UserErrors `json:"userErrors"`
		} `json:"customerEmailMarketingConsentUpdate" graphql:"customerEmailMarketingConsentUpdate(input: $input)"`
	}
	if err := r.client.Mutation(ctx, &query, variables); err != nil {
		return nil, err
	}

	for _, v := range query.UpdateEmailMarketing.UserErrors {
		return nil, errors.New(string(v.Message))
	}

	return &query.UpdateEmailMarketing.Customer, nil
}
