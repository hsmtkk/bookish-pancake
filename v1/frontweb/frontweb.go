package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"

	"github.com/hsmtkk/bookish-pancake/proto"
	"github.com/hsmtkk/bookish-pancake/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	backURL, err := util.RequiredEnvVar("BACK_URL")
	if err != nil {
		log.Fatal(err)
	}
	backHost, err := util.GetHostFromURL(backURL)
	if err != nil {
		log.Fatal(err)
	}
	port, err := util.GetPort()
	if err != nil {
		log.Fatal(err)
	}

	hdl, err := newHandler(backHost)
	if err != nil {
		log.Fatal(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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
	grpcConn, err := grpc.Dial(fmt.Sprintf("%s:443", backHost), grpc.WithTransportCredentials(cred))
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
	city := ectx.FormValue("city")
	resp, err := h.grpcClient.GetWeather(ectx.Request().Context(), &proto.WeatherRequest{City: city})
	if err != nil {
		return fmt.Errorf("proto.WeatherServiceClient.GetWeather failed; %w", err)
	}
	temp := resp.GetTemperature()
	pressure := resp.GetPressure()
	humidity := resp.GetHumidity()
	html := fmt.Sprintf(postHTML, temp, pressure, humidity)
	return ectx.HTML(http.StatusOK, html)
}
