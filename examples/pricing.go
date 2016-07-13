package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/ZhongtianWang/azure/ratecard"
)

const (
	credentialsPath = "/.azure/credentials.json"
)

type credentials struct {
	ClientID       string `json:"clientID"`
	ClientSecret   string `json:"clientSecret"`
	TenantID       string `json:"tenantID"`
	SubscriptionID string `json:"subscriptionID"`
}

func withInspection() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			fmt.Printf("Inspecting Request: %s %s\n", r.Method, r.URL)
			return p.Prepare(r)
		})
	}
}

func byInspecting() autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			fmt.Printf("Inspecting Response: %s for %s %s\n", resp.Status, resp.Request.Method, resp.Request.URL)
			return r.Respond(resp)
		})
	}
}

func loadCredentials(cred *credentials) error {
	u, err := user.Current()
	if err != nil {
		return errors.New("unable to determine current user")
	}

	dir := u.HomeDir + credentialsPath
	f, err := os.Open(dir)
	if err != nil {
		return errors.New("unable to open Azure credentials at " + dir)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.New("unable to read " + dir)
	}

	if err = json.Unmarshal(b, cred); err != nil {
		return errors.New(dir + " contains invalid JSON")
	}

	return nil
}

func main() {
	cred := credentials{}
	if err := loadCredentials(&cred); err != nil {
		log.Fatalf("Error: %v", err)
	}

	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(cred.TenantID)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	spt, err := azure.NewServicePrincipalToken(*oauthConfig, cred.ClientID,
		cred.ClientSecret, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	client := ratecard.NewClient(cred.SubscriptionID)
	client.Authorizer = spt
	client.RequestInspector = withInspection()
	client.ResponseInspector = byInspecting()

	param := ratecard.RateCardGetParameters{
		OfferDurableId: stringPtr("MS-AZR-0003p"),
		Currency:       stringPtr("USD"),
		Locale:         stringPtr("en-US"),
		RegionInfo:     stringPtr("US"),
	}

	rateCard, err := client.Get(param, make(chan struct{}))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Meters:")
	for _, meter := range *rateCard.Meters {
		log.Println("MeterId: ", *meter.MeterId)
		log.Println("MeterCategory: ", *meter.MeterCategory)
		log.Println("MeterSubCategory: ", *meter.MeterSubCategory)
		log.Println("Unit: ", *meter.Unit)
		log.Println("MeterRates:")
		for quantity, rate := range *meter.MeterRates {
			log.Println(quantity, ":", *rate)
		}
		log.Println("EffectiveDate: ", *meter.EffectiveDate)
		log.Println("IncludedQuantity: ", *meter.IncludedQuantity)
	}
}

func stringPtr(s string) *string {
	return &s
}
