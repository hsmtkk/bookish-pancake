package openweather

import (
	"context"
	"net/http"

	"github.com/hsmtkk/bookish-pancake/openweather"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func GetCurrentWeather(ctx context.Context, apiKey, city string) (openweather.WeatherData, error) {
	client := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	return openweather.New(client).GetCurrentWeather(ctx, apiKey, city)
}
