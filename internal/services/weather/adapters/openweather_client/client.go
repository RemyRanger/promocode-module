package openweather_client

import (
	"APIs/internal/common/config"
	"APIs/internal/common/entities"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/rs/zerolog/log"
)

// timeoutClient : timeout in second
const timeoutClient = 20

// APIClient : client struct
type OpenWeatherClient struct {
	Client *ClientWithResponses
	Appid  string
}

// NewClientAPI : initializes a new ClientYanport struct.
func NewClientAPI(app_config config.Config) *OpenWeatherClient {
	client, err := NewClient(app_config.Openweather.Url)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create OpenWeather client.")
	}

	client.Client = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
		Timeout:   timeoutClient * time.Second,
	}

	return &OpenWeatherClient{
		Client: &ClientWithResponses{client},
		Appid:  app_config.Openweather.Apikey,
	}
}

func (c *OpenWeatherClient) FetchWeather(ctx context.Context, town string) (*entities.Weather, error) {
	metric := FetchWeatherParamsUnits("metric") // metric = Set Units of measurement to Celsius
	params := FetchWeatherParams{
		Q:     &town,
		Units: &metric,
		Appid: c.Appid,
	}
	weatherResponse, err := c.Client.FetchWeatherWithResponse(ctx, &params)
	if err != nil {
		return nil, err
	}

	switch weatherResponse.StatusCode() {
	case http.StatusOK:
		weathers := *weatherResponse.JSON200.Weather
		mainType := string(*weathers[0].Main)

		result := &entities.Weather{
			Town: town,
			Temp: *weatherResponse.JSON200.Main.Temp,
			Type: entities.WeatherType(strings.ToLower(mainType)),
		}
		return result, nil
	case http.StatusBadRequest:
		return nil, errors.New("receive bad request from OpenWeather API")
	default:
		return nil, errors.New("receive error from OpenWeather API")
	}
}
