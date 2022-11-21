package openweather

import (
	"context"
	"net/http"

	"github.com/hsmtkk/bookish-pancake/openweather"
)

func GetCurrentWeather(ctx context.Context, apiKey, city string) (openweather.WeatherData, error) {
	return openweather.New(http.DefaultClient).GetCurrentWeather(ctx, apiKey, city)
}
