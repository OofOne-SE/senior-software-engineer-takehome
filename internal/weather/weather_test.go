package weather_test

import (
	"senior-software-engineer-takehome/client"
	"senior-software-engineer-takehome/client/weather_assets"
	"senior-software-engineer-takehome/models"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

func TestAddWeatherUnit(t *testing.T) {
	clt := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:     "localhost:8080",
		BasePath: client.DefaultBasePath,
		Schemes:  []string{"http"},
	})

	postRes, err := clt.WeatherAssets.PostDataWeather(&weather_assets.PostDataWeatherParams{
		Body: &models.WeatherUnit{
			Date:        (*strfmt.Date)(swag.Time(time.Now())),
			Humidity:    swag.Float64(42.42),
			Temperature: swag.Float64(42.42),
		},
	})
	assert.NotNil(t, postRes)
	assert.Nil(t, err)
}

func TestGetWeatherUnits(t *testing.T) {
	clt := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Host:     "localhost:8080",
		BasePath: client.DefaultBasePath,
		Schemes:  []string{"http"},
	})

	getRes, err := clt.WeatherAssets.GetDataWeather(&weather_assets.GetDataWeatherParams{
		From: strfmt.Date(time.Unix(1674918395, 0)),
		To:   strfmt.Date(time.Unix(1674918395, 0)),
	})
	assert.NotNil(t, getRes)
	assert.Empty(t, getRes.Payload)
	assert.Nil(t, err)
}
