package openweather

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"github.com/hsmtkk/bookish-pancake/city"
	"github.com/hsmtkk/bookish-pancake/constant"
)

func GetCityData(ctx context.Context, apiKey, cityName string) ([]city.CityData, error) {
	tracer := otel.GetTracerProvider().Tracer(constant.ServiceName)
	spanCtx, span := tracer.Start(ctx, "GetCityData")
	defer span.End()

	clt := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	return city.New(clt).GetCitiesData(spanCtx, apiKey, cityName)
}
