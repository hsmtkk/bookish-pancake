package city

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hsmtkk/bookish-pancake/util"
)

type CityData struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CitiesDataGetter interface {
	GetCitiesData(ctx context.Context, apiKey, city string) ([]CityData, error)
}

type impl struct {
	clt     *http.Client
	baseURL string
}

func New(clt *http.Client) CitiesDataGetter {
	baseURL := "https://api.api-ninjas.com"
	return &impl{clt, baseURL}
}

func NewForTest(clt *http.Client, baseURL string) CitiesDataGetter {
	return &impl{clt, baseURL}
}

func (i *impl) GetCitiesData(ctx context.Context, apiKey, city string) ([]CityData, error) {
	url := fmt.Sprintf("%s/v1/city/?name=%s", i.baseURL, city)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext failed; %w", err)
	}
	req.Header.Add("X-Api-Key", apiKey)
	resp, err := i.clt.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do failed; %w", err)
	}
	defer resp.Body.Close()
	bs, err := util.HandleHTTPResponse(resp)
	if err != nil {
		return nil, err
	}
	var cities []CityData
	if err := json.Unmarshal(bs, &cities); err != nil {
		return nil, fmt.Errorf("json.Unmarshal failed; %w", err)
	}
	return cities, nil
}
