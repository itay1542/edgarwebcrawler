package requests

import (
	"encoding/json"
	"fmt"
	"github.com/itay1542/edgarwebcrawler/DAL"
	"io"
	"net/http"
	"time"
)

type CompanyGetter interface {
	GetCompanyDetails(symbol string) (*CompanyDetails, error)
}

type CompanyDetails struct {
	Sector    string            `json:"Sector"`
	Exchange  DAL.StockExchange `json:"Exchange"`
	Name      string            `json:"Name"`
	MarketCap string            `json:"MarketCapitalization"`
}

//AlphaVantageRequester implements CompanyGetter
type AlphaVantageRequester struct {
	host, apiKey string
}

func NewAlphaVantageRequester(host, apiKey string) *AlphaVantageRequester {
	return &AlphaVantageRequester{
		host:   host,
		apiKey: apiKey,
	}
}

func (a AlphaVantageRequester) GetCompanyDetails(symbol string) (*CompanyDetails, error) {
	url := fmt.Sprintf("%s/query?function=OVERVIEW&symbol=%s&apikey=%s", a.host, symbol, a.apiKey)
	t := &http.Transport{
		TLSHandshakeTimeout: 2 * time.Second,
	}
	c := &http.Client{
		Transport: t,
	}
	response, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("received response status code %d", response.StatusCode)
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var companyDetails CompanyDetails
	err = json.Unmarshal(data, &companyDetails)
	if err != nil {
		return nil, err
	}
	return &companyDetails, nil
}
