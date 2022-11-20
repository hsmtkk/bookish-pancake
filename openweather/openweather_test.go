package openweather_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hsmtkk/bookish-pancake/openweather"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentWeather(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokyo, err := os.ReadFile("tokyo.json")
		assert.Nil(t, err)
		io.Copy(w, bytes.NewReader(tokyo))
	}))
	defer ts.Close()

	want := openweather.WeatherData{
		Main: openweather.WeatherDataMain{
			Temperature: 284.5,
			Pressure:    1025,
			Humidity:    69,
		},
	}
	got, err := openweather.NewForTest(ts.Client(), ts.URL).GetCurrentWeather(context.Background(), "secret", "Tokyo")
	assert.Equal(t, want, got)
	assert.Nil(t, err)
}
