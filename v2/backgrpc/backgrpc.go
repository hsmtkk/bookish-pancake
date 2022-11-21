package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hsmtkk/bookish-pancake/constant"
	"github.com/hsmtkk/bookish-pancake/proto"
	"github.com/hsmtkk/bookish-pancake/utilenv"
	"github.com/hsmtkk/bookish-pancake/utilgcp"
	"github.com/hsmtkk/bookish-pancake/v2/openweather"
	"github.com/hsmtkk/bookish-pancake/v2/provider"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer(constant.ServiceName)
}

type server struct {
	proto.UnimplementedWeatherServiceServer
	apiKey string
}

func newServer(apiKey string) *server {
	return &server{apiKey: apiKey}
}

func (s *server) GetWeather(ctx context.Context, in *proto.WeatherRequest) (*proto.WeatherResponse, error) {
	spanCtx, span := tracer.Start(ctx, "GetWeather")
	defer span.End()

	data, err := openweather.GetCurrentWeather(spanCtx, s.apiKey, in.GetCity())
	if err != nil {
		return nil, err
	}
	return &proto.WeatherResponse{
		Temperature: data.Main.Temperature,
		Pressure:    int64(data.Main.Pressure),
		Humidity:    int64(data.Main.Humidity),
	}, nil
}

func main() {
	ctx := context.Background()
	prov, err := provider.Provider(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer prov.ForceFlush(ctx)
	otel.SetTracerProvider(prov)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	apiKey, err := utilgcp.GetSecret(context.Background(), "openweather_api_key")
	if err != nil {
		log.Fatal(err)
	}
	port, err := utilenv.GetPort()
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("net.Listen failed; %v", err.Error())
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()), grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()))
	proto.RegisterWeatherServiceServer(s, newServer(apiKey))
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
