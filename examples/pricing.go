// This example lists all Linux Pay-as-you-go VM Pricing and Sizes in USD for all US VMs.

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

	"github.com/Azure/azure-sdk-for-go/arm/compute"
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

type rate struct {
	size   string
	price  float32
	region string
	cpu    int32
	ram    float32
	disk   int32
}

// id is the unique identifier for VMs..
type id struct {
	region, size string
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

	rateCardClient := ratecard.NewClient(cred.SubscriptionID)
	rateCardClient.Authorizer = spt

	vmSizeClient := compute.NewVirtualMachineSizesClient(cred.SubscriptionID)
	vmSizeClient.Authorizer = spt

	// Uncomment to inspect http request and responses
	// rateCardClient.RequestInspector = withInspection()
	// rateCardClient.ResponseInspector = byInspecting()
	// vmSizeClient.RequestInspector = withInspection()
	// vmSizeClient.ResponseInspector = byInspecting()

	// List all Linux only Pay-as-you-go VM Pricing in USD for all US VMs.
	// Frontend to backend linux VM sizes mapping.
	linuxVms := map[string]string{
		"BASIC.A0":        "Basic_A0",
		"BASIC.A1":        "Basic_A1",
		"BASIC.A2":        "Basic_A2",
		"BASIC.A3":        "Basic_A3",
		"BASIC.A4":        "Basic_A4",
		"A0":              "Standard_A0",
		"A1":              "Standard_A1",
		"A2":              "Standard_A2",
		"A3":              "Standard_A3",
		"A4":              "Standard_A4",
		"A5":              "Standard_A5",
		"A6":              "Standard_A6",
		"A7":              "Standard_A7",
		"A8":              "Standard_A8",
		"A9":              "Standard_A9",
		"Standard_D1":     "Standard_D1",
		"Standard_D2":     "Standard_D2",
		"Standard_D3":     "Standard_D3",
		"Standard_D4":     "Standard_D4",
		"Standard_D11":    "Standard_D11",
		"Standard_D12":    "Standard_D12",
		"Standard_D13":    "Standard_D13",
		"Standard_D14":    "Standard_D14",
		"Standard_D1_v2":  "Standard_D1_v2",
		"Standard_D2_v2":  "Standard_D2_v2",
		"Standard_D3_v2":  "Standard_D3_v2",
		"Standard_D4_v2":  "Standard_D4_v2",
		"Standard_D5_v2":  "Standard_D5_v2",
		"Standard_D11_v2": "Standard_D11_v2",
		"Standard_D12_v2": "Standard_D12_v2",
		"Standard_D13_v2": "Standard_D13_v2",
		"Standard_D14_v2": "Standard_D14_v2",
		"Standard_D15_v2": "Standard_D15_v2",
		"Standard_DS1":    "Standard_DS1",
		"Standard_DS2":    "Standard_DS2",
		"Standard_DS3":    "Standard_DS3",
		"Standard_DS4":    "Standard_DS4",
		"Standard_DS11":   "Standard_DS11",
		"Standard_DS12":   "Standard_DS12",
		"Standard_DS13":   "Standard_DS13",
		"Standard_DS14":   "Standard_DS14",
		"Standard_G1":     "Standard_G1",
		"Standard_G2":     "Standard_G2",
		"Standard_G3":     "Standard_G3",
		"Standard_G4":     "Standard_G4",
		"Standard_G5":     "Standard_G5",
		"Standard_GS1":    "Standard_GS1",
		"Standard_GS2":    "Standard_GS2",
		"Standard_GS3":    "Standard_GS3",
		"Standard_GS4":    "Standard_GS4",
		"Standard_GS5":    "Standard_GS5",
		"Standard_F1":     "Standard_F1",
		"Standard_F2":     "Standard_F2",
		"Standard_F4":     "Standard_F4",
		"Standard_F8":     "Standard_F8",
		"Standard_F16":    "Standard_F16",
	}

	// Frontend to backend US locations mapping.
	usLocations := map[string]string{
		"US East":          "eastus",
		"US East 2":        "eastus2",
		"US West":          "westus",
		"US Central":       "centralus",
		"US North Central": "northcentralus",
		"US South Central": "southcentralus",
		"US West 2":        "westus2",
		"US West Central":  "westcentralus",
	}

	// The current RateCard API returns weird response. To filter out only Linux VMs, We
	// create a map that maps the API formatted size tyoes to the actual ones.
	filterSet := make(map[string]string)
	for frontSize, backSize := range linuxVms {
		meterSize := fmt.Sprintf("%s VM", frontSize)
		filterSet[meterSize] = backSize
	}

	param := ratecard.RateCardGetParameters{
		OfferDurableId: stringPtr("MS-AZR-0003p"),
		Currency:       stringPtr("USD"),
		Locale:         stringPtr("en-US"),
		RegionInfo:     stringPtr("US"),
	}

	rateCard, err := rateCardClient.Get(param, make(chan struct{}))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	result := make(map[id]rate)

	for _, meter := range *rateCard.Meters {
		if *meter.MeterCategory != ratecard.VirtualMachines {
			continue
		}

		meterSize := *meter.MeterSubCategory
		price := *(*meter.MeterRates)["0"]
		region := *meter.MeterRegion

		// Filter out Linux VMs.
		backSize, ok := filterSet[meterSize]
		if !ok {
			continue
		}

		// Filter out US regions.
		backRegion, ok := usLocations[region]
		if !ok {
			continue
		}

		result[id{region: backRegion, size: backSize}] = rate{
			size:   backSize,
			price:  price,
			region: backRegion,
		}
	}

	for _, backRegion := range usLocations {
		vmSizes, err := vmSizeClient.List(backRegion)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		for _, vmSize := range *vmSizes.Value {
			backSize := *vmSize.Name
			key := id{region: backRegion, size: backSize}
			rate, ok := result[key]
			if !ok {
				// log.Printf("%s:%s Found in vm info but not in ratecard", backRegion, backSize)
				continue
			}

			rate.cpu = *vmSize.NumberOfCores
			rate.disk = *vmSize.ResourceDiskSizeInMB / 1024
			rate.ram = (float32)(*vmSize.MemoryInMB) / 1024

			result[key] = rate
		}
	}

	for _, rate := range result {
		log.Printf("%v.\n", rate)
	}

}

func stringPtr(s string) *string {
	return &s
}
