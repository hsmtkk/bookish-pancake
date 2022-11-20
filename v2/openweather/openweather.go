package openweather

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"github.com/hsmtkk/bookish-pancake/constant"
	"github.com/hsmtkk/bookish-pancake/openweather"
)

func GetCurrentWeather(ctx context.Context, apiKey, city string) (openweather.WeatherData, error) {
	tracer := otel.GetTracerProvider().Tracer(constant.ServiceName)
	spanCtx, span := tracer.Start(ctx, "GetCurrentWeather")
	defer span.End()

	clt := &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	return openweather.New(clt).GetCurrentWeather(spanCtx, apiKey, city)
}
