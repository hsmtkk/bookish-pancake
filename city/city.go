package city

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hsmtkk/bookish-pancake/util"
)

type CityData struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type CityDataGetter interface {
	GetCityData(ctx context.Context, apiKey, city string) (CityData, error)
}

type impl struct {
	clt     *http.Client
	baseURL string
}

func New(clt *http.Client) CityDataGetter {
	baseURL := "https://api.api-ninjas.com"
	return &impl{clt, baseURL}
}

func NewForTest(clt *http.Client, baseURL string) CityDataGetter {
	return &impl{clt, baseURL}
}

func (i *impl) GetCityData(ctx context.Context, apiKey, city string) (CityData, error) {
	url := fmt.Sprintf("%s/v1/city/?name=%s", i.baseURL, city)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return CityData{}, fmt.Errorf("http.NewRequestWithContext failed; %w", err)
	}
	req.Header.Add("X-Api-Key", apiKey)
	resp, err := i.clt.Do(req)
	if err != nil {
		return CityData{}, fmt.Errorf("http.Client.Do failed; %w", err)
	}
	defer resp.Body.Close()
	bs, err := util.HandleHTTPResponse(resp)
	if err != nil {
		return CityData{}, err
	}
	log.Printf("%s\n", string(bs))
	var encoded []CityData
	if err := json.Unmarshal(bs, &encoded); err != nil {
		return CityData{}, fmt.Errorf("json.Unmarshal failed; %w", err)
	}
	if len(encoded) == 1 {
		return encoded[0], nil
	} else {
		return CityData{}, fmt.Errorf("failed to get city data")
	}
}
