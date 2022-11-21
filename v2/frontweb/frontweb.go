package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/hsmtkk/bookish-pancake/constant"
	"github.com/hsmtkk/bookish-pancake/proto"
	"github.com/hsmtkk/bookish-pancake/utilenv"
	"github.com/hsmtkk/bookish-pancake/v2/profiler"
	"github.com/hsmtkk/bookish-pancake/v2/provider"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer(constant.ServiceName)
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

	backURL, err := utilenv.RequiredVar("BACK_URL")
	if err != nil {
		log.Fatal(err)
	}
	parsed, err := url.Parse(backURL)
	if err != nil {
		log.Fatalf("url.Parse failed; %v", err.Error())
	}
	port, err := utilenv.GetPort()
	if err != nil {
		log.Fatal(err)
	}

	hdl, err := newHandler(parsed.Host)
	if err != nil {
		log.Fatal(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware(constant.ServiceName))

	// Routes
	e.GET("/", hdl.getIndex)
	e.POST("/", hdl.postIndex)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type handler struct {
	grpcClient proto.WeatherServiceClient
}

func newHandler(backHost string) (*handler, error) {
	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("x509.SystemCertPool failed; %w", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		RootCAs: systemRoots,
	})
	grpcConn, err := grpc.Dial(fmt.Sprintf("%s:443", backHost), grpc.WithTransportCredentials(cred), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()), grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	if err != nil {
		return nil, fmt.Errorf("grpc.Dial failed; %w", err)
	}
	grpcClient := proto.NewWeatherServiceClient(grpcConn)
	return &handler{grpcClient}, nil
}

const indexHTML = `<html>
 <body>
  <form method="post" action="/">
   <input type="text" name="city">
   <input type="submit" value="submit">
  </form>
 </body>
</html>
`

func (h *handler) getIndex(ectx echo.Context) error {
	_, span := tracer.Start(ectx.Request().Context(), "getIndex")
	defer span.End()

	return ectx.HTML(http.StatusOK, indexHTML)
}

const postHTML = `<html>
 <body>
  <ul>
   <li>Temperature: %f</li>
   <li>Pressure: %d</li>
   <li>Humidityd: %d</li>
  </ul>
 </body>
</html>`

func (h *handler) postIndex(ectx echo.Context) error {
	spanCtx, span := tracer.Start(ectx.Request().Context(), "getIndex")
	defer span.End()

	city := ectx.FormValue("city")
	resp, err := h.grpcClient.GetWeather(spanCtx, &proto.WeatherRequest{City: city})
	if err != nil {
		return fmt.Errorf("proto.WeatherServiceClient.GetWeather failed; %w", err)
	}
	temp := resp.GetTemperature()
	pressure := resp.GetPressure()
	humidity := resp.GetHumidity()
	html := fmt.Sprintf(postHTML, temp, pressure, humidity)
	return ectx.HTML(http.StatusOK, html)
}
