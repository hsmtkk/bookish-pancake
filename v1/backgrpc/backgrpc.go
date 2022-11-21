package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hsmtkk/bookish-pancake/proto"
	"github.com/hsmtkk/bookish-pancake/util"
	"github.com/hsmtkk/bookish-pancake/v1/openweather"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedWeatherServiceServer
	apiKey string
}

func newServer(apiKey string) *server {
	return &server{apiKey: apiKey}
}

func (s *server) GetWeather(ctx context.Context, in *proto.WeatherRequest) (*proto.WeatherResponse, error) {
	data, err := openweather.GetCurrentWeather(ctx, s.apiKey, in.GetCity())
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
	apiKey, err := util.GetSecret(context.Background(), "openweather_api_key")
	if err != nil {
		log.Fatal(err)
	}
	port, err := util.GetPort()
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("net.Listen failed; %v", err.Error())
	}
	s := grpc.NewServer()
	proto.RegisterWeatherServiceServer(s, newServer(apiKey))
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
