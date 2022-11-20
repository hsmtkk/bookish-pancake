package openweather

import (
	"context"
	"net/http"

	"github.com/hsmtkk/bookish-pancake/city"
)

func GetCityData(ctx context.Context, apiKey, cityName string) ([]city.CityData, error) {
	return city.New(http.DefaultClient).GetCitiesData(context.Background(), apiKey, cityName)
}
