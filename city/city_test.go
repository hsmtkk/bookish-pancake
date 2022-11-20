package city_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hsmtkk/bookish-pancake/city"
	"github.com/stretchr/testify/assert"
)

func TestGetCityData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokyo, err := os.ReadFile("tokyo.json")
		assert.Nil(t, err)
		io.Copy(w, bytes.NewReader(tokyo))
	}))
	defer ts.Close()

	want := city.CityData{
		Name:      "Tokyo",
		Latitude:  35.6897,
		Longitude: 139.692,
	}
	got, err := city.NewForTest(ts.Client(), ts.URL).GetCityData(context.Background(), "secret", "Tokyo")
	assert.Equal(t, want, got)
	assert.Nil(t, err)
}
