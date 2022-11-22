package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/hsmtkk/bookish-pancake/constant"
	"github.com/hsmtkk/bookish-pancake/proto"
	"github.com/hsmtkk/bookish-pancake/utilenv"
	"github.com/hsmtkk/bookish-pancake/utilgcp"
	"github.com/hsmtkk/bookish-pancake/v2/openweather"
	"github.com/hsmtkk/bookish-pancake/v2/profiler"
	"github.com/hsmtkk/bookish-pancake/v2/provider"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/genproto/googleapis/api/metric"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	before := time.Now()
	data, err := openweather.GetCurrentWeather(spanCtx, s.apiKey, in.GetCity())
	if err != nil {
		return nil, err
	}
	elapsed := time.Since(before)
	if err := s.recordWebLatency(spanCtx, elapsed.Milliseconds()); err != nil {
		log.Printf("error: %v", err.Error()) // accept error
	}

	return &proto.WeatherResponse{
		Temperature: data.Main.Temperature,
		Pressure:    int64(data.Main.Pressure),
		Humidity:    int64(data.Main.Humidity),
	}, nil
}

const metricType = "custom.googleapis.com/weblatency"

func (s *server) recordWebLatency(ctx context.Context, latencyMilliSeconds int64) error {
	spanCtx, span := tracer.Start(ctx, "recordWebLatency")
	defer span.End()

	projectID, err := utilgcp.ProjectID(spanCtx)
	if err != nil {
		return err
	}
	clt, err := monitoring.NewMetricClient(spanCtx)
	if err != nil {
		return fmt.Errorf("monitoring.NewMetricClient failed; %w", err)
	}
	defer clt.Close()
	now := timestamppb.Now()
	req := monitoringpb.CreateTimeSeriesRequest{
		Name: "projects/" + projectID,
		TimeSeries: []*monitoringpb.TimeSeries{{
			Metric: &metric.Metric{
				Type: metricType,
			},
			Points: []*monitoringpb.Point{{
				Interval: &monitoringpb.TimeInterval{
					StartTime: now,
					EndTime:   now,
				},
				Value: &monitoringpb.TypedValue{
					Value: &monitoringpb.TypedValue_Int64Value{
						Int64Value: latencyMilliSeconds,
					},
				},
			}},
		}},
	}
	if err := clt.CreateTimeSeries(spanCtx, &req); err != nil {
		return fmt.Errorf("monitoring.MetricClient.CreateTimeSeries failed; %w", err)
	}
	return nil
}

func main() {
	ctx := context.Background()
	if err := profiler.Start(ctx); err != nil {
		log.Fatal(err)
	}
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
