package ratecard

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

const (
	// APIVersion is required for RateCard api.
	APIVersion = "2015-06-01-preview"

	// DefaultBaseURI is the default url for RateCard api.
	DefaultBaseURI = "https://management.azure.com"
)

// Client is the base Azure RateCard client.
type Client struct {
	autorest.Client
	BaseURI        string
	SubscriptionID string
}

// NewClient creates an instance of RateCard Client
func NewClient(subscriptionID string) Client {
	return Client{
		Client:         autorest.NewClientWithUserAgent(""),
		BaseURI:        DefaultBaseURI,
		SubscriptionID: subscriptionID,
	}
}

func (client Client) Get(parameters RateCardGetParameters, cancel <-chan struct{}) (result RateCard, err error) {
	req, err := client.GetPreparer(parameters, cancel)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "RateCardClient", "Get", nil, "Failure preparing request")
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "RateCardClient", "Get", resp, "Failure sending request")
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "RateCardClient", "Get", resp, "Failure responding to request")
	}

	return
}

func (client Client) GetPreparer(parameters RateCardGetParameters, cancel <-chan struct{}) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"subscriptionId": url.QueryEscape(client.SubscriptionID),
	}

	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
		"$filter": fmt.Sprintf("OfferDurableId eq '%s' and Currency eq '%s' and Locale eq '%s' and RegionInfo eq '%s'",
			*parameters.OfferDurableId, *parameters.Currency, *parameters.Locale, *parameters.RegionInfo),
	}

	return autorest.Prepare(&http.Request{Cancel: cancel},
		autorest.AsJSON(),
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPath("/subscriptions/{subscriptionId}/providers/Microsoft.Commerce/RateCard"),
		autorest.WithPathParameters(pathParameters),
		autorest.WithQueryParameters(queryParameters))
}

func (client Client) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

func (client Client) GetResponder(resp *http.Response) (result RateCard, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
