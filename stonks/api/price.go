package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	BaseUrl         string = "apidojo-yahoo-finance-v1.p.rapidapi.com/market/get-quotes"
	HostHeader      string = "apidojo-yahoo-finance-v1.p.rapidapi.com"
	TimeoutDuration string = "5s"
	ErrNoneFound    error  = errors.New("No stock price data found")
)

type PriceApiClient struct {
	apiKey string
	client http.Client
}

func NewPriceApiClient(apiKey string) (*PriceApiClient, error) {

	fmtDur, err := time.ParseDuration(TimeoutDuration)
	if err != nil {
		return nil, err
	}

	c := http.Client{
		Timeout: fmtDur,
	}

	url := formatHealthCheckUrl(BaseUrl)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-rapidapi-host", HostHeader)
	req.Header.Add("x-rapidapi-key", apiKey)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to call price API: %v", resp.StatusCode)
	}

	pc := &PriceApiClient{
		apiKey: apiKey,
		client: c,
	}

	return pc, nil
}

func (p *PriceApiClient) GetPrices(symbols []string) ([]Stonk, error) {
	var output []Stonk

	if len(symbols) < 1 {
		return output, errors.New("must provide at least 1 stock symbol")
	}

	url := formatPriceUrl(BaseUrl, symbols)

	// make request bb
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return output, err
	}
	req.Header.Add("x-rapidapi-host", HostHeader)
	req.Header.Add("x-rapidapi-key", p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return output, err
	}

	if resp.StatusCode != 200 {
		return output, fmt.Errorf("Invalid request: %v", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Unable to read price API body: %v", err)
		return output, err
	}

	var pricePayload PriceApiResponse
	err = json.Unmarshal(respBytes, &pricePayload)
	if err != nil {
		log.Printf("Unable to unmarshal response: %v", err)
		return output, err
	}
	defer resp.Body.Close()

	if len(pricePayload.QuoteResponse.Result) < 1 {
		return output, ErrNoneFound
	}

	output = pricePayload.QuoteResponse.Result

	return output, nil
}

func formatHealthCheckUrl(baseUrl string) string {

	return fmt.Sprintf("https://%v?symbols=SPY", baseUrl)
}

func formatPriceUrl(baseUrl string, symbols []string) string {

	var symbolParam string
	for _, s := range symbols {
		symbolParam += fmt.Sprintf("%v,", s)
	}

	symbolParam = strings.TrimSuffix(symbolParam, ",")
	return fmt.Sprintf("https://%v?symbols=%v", baseUrl, symbolParam)

}
