package input

import "time"

// CustomerInput represents the input object for creating or updating a customer.
type CustomerInput struct {
	Id                    *string                        `json:"id,omitempty"`
	Email                 *string                        `json:"email,omitempty"`
	Phone                 *string                        `json:"phone,omitempty"`
	FirstName             *string                        `json:"firstName,omitempty"`
	LastName              *string                        `json:"lastName,omitempty"`
	Locale                *string                        `json:"locale,omitempty"`
	Note                  *string                        `json:"note,omitempty"`
	TaxExempt             *bool                          `json:"taxExempt,omitempty"`
	TaxExemptions         []string                       `json:"taxExemptions,omitempty"`
	Tags                  []string                       `json:"tags,omitempty"`
	SmsMarketingConsent   *CustomerMarketingConsentInput `json:"smsMarketingConsent,omitempty"`
	EmailMarketingConsent *CustomerMarketingConsentInput `json:"emailMarketingConsent,omitempty"`
	Addresses             []MailingAddressInput          `json:"addresses,omitempty"`
	Metafields            []MetafieldInput               `json:"metafields,omitempty"`
}

type CustomerMarketingConsentInput struct {
	ConsentUpdatedAt    *time.Time `json:"consentUpdatedAt,omitempty"`
	MarketingOptInLevel *string    `json:"marketingOptInLevel,omitempty"`
	MarketingState      *string    `json:"marketingState,omitempty"`
}

type MailingAddressInput struct {
	Address1     *string `json:"address1,omitempty"`
	Address2     *string `json:"address2,omitempty"`
	City         *string `json:"city,omitempty"`
	Company      *string `json:"company,omitempty"`
	CountryCode  *string `json:"countryCode,omitempty"`
	FirstName    *string `json:"firstName,omitempty"`
	LastName     *string `json:"lastName,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	ProvinceCode *string `json:"provinceCode,omitempty"`
	Zip          *string `json:"zip,omitempty"`
}
