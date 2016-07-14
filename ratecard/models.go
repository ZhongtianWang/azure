package ratecard

import "github.com/Azure/go-autorest/autorest"

type MeterCategory string

const (
	VirtualMachines MeterCategory = "Virtual Machines"
	CloudServices   MeterCategory = "Cloud Services"
	Networking      MeterCategory = "Networking"
	Storage         MeterCategory = "Storage"
	DataServices    MeterCategory = "Data Services"
)

type RateCard struct {
	autorest.Response `json:"-"`
	OfferTerms        *[]interface{} `json:"OfferTerms,omitempty"`
	Meters            *[]Meters      `json:"Meters,omitempty"`
	Currency          *string        `json:"Currency,omitempty"`
	Locale            *string        `json:"Locale,omitempty"`
	IsTaxIncluded     *bool          `json:"IsTaxIncluded,omitempty"`
	Tags              *[]string      `json:"Tags,omitempty"`
}

type Meters struct {
	autorest.Response `json:"-"`
	MeterId           *string              `json:"MeterId,omitempty"`
	MeterCategory     *MeterCategory       `json:"MeterCategory,omitempty"`
	MeterSubCategory  *string              `json:"MeterSubCategory,omitempty"`
	Unit              *string              `json:"Unit,omitempty"`
	MeterRegion       *string              `json:"MeterRegion,omitempty"`
	MeterRates        *map[string]*float32 `json:"MeterRates,omitempty"`
	EffectiveDate     *string              `json:"EffectiveDate,omitempty"`
	IncludedQuantity  *float32             `json:"IncludedQuantity,omitempty"`
}

type RateCardGetParameters struct {
	OfferDurableId *string
	Currency       *string
	Locale         *string
	RegionInfo     *string
}
