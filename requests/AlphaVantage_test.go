package requests

import (
	edgarwebcrawler "github.com/itay1542/edgarwebcrawler/DAL"
	"testing"
)

func TestAlphaVantageRequester_GetCompanyDetails(t *testing.T) {
	requester := &AlphaVantageRequester{AlphaVantageConfig{
		host:   "https://www.alphavantage.co",
		apiKey: "demo",
	}}
	details, error := requester.GetCompanyDetails("IBM")
	if error != nil {
		t.Fatalf("failed to request, %s", error)
	}
	expectedName := "International Business Machines Corporation"
	gotName := details.Name
	if gotName != expectedName {
		t.Fatalf("Received %s, expected %s", gotName, details.Name)
	}
	expectedExchange := edgarwebcrawler.NYSE
	gotExchange := details.Exchange
	if gotExchange != expectedExchange {
		t.Fatalf("Received %s, expected %s", gotExchange, expectedExchange)
	}

}
