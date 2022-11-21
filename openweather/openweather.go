package openweather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hsmtkk/bookish-pancake/util"
)

type WeatherDataMain struct {
	Temperature float64 `json:"temp"`
	Pressure    int     `json:"pressure"`
	Humidity    int     `json:"humidity"`
}

type WeatherData struct {
	Main WeatherDataMain `json:"main"`
}

type CurrentWeatherGetter interface {
	GetCurrentWeather(ctx context.Context, apiKey, city string) (WeatherData, error)
}

type impl struct {
	clt     *http.Client
	baseURL string
}

func New(clt *http.Client) CurrentWeatherGetter {
	baseURL := "https://api.openweathermap.org"
	return &impl{clt, baseURL}
}

func NewForTest(clt *http.Client, baseURL string) CurrentWeatherGetter {
	return &impl{clt, baseURL}
}

func (i *impl) GetCurrentWeather(ctx context.Context, apiKey, city string) (WeatherData, error) {
	url := fmt.Sprintf("%s/data/2.5/weather?q=%s&appid=%s", i.baseURL, city, apiKey)
	var result WeatherData
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, fmt.Errorf("http.NewRequestWithContext failed; %w", err)
	}
	resp, err := i.clt.Do(req)
	if err != nil {
		return result, fmt.Errorf("http.Client.Do failed; %w", err)
	}
	defer resp.Body.Close()
	bs, err := util.HandleHTTPResponse(resp)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(bs, &result); err != nil {
		return result, fmt.Errorf("json.Unmarshal failed; %w", err)
	}
	return result, nil
}
