package test_utils

import (
	"APIs/internal/services/weather/adapters/openweather_client"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/onsi/gomega/ghttp"

	. "github.com/onsi/gomega"
)

const respBodyLyon = `{
	"coord": {
		"lon": 4.5833,
		"lat": 45.75
	},
	"weather": [
		{
		"id": 800,
		"main": "Clear",
		"description": "clear sky",
		"icon": "01d"
		}
	],
	"base": "stations",
	"main": {
		"temp": 21.49,
		"feels_like": 21.1,
		"temp_min": 20.56,
		"temp_max": 22.33,
		"pressure": 1010,
		"humidity": 54,
		"sea_level": 1010,
		"grnd_level": 950
	},
	"visibility": 10000,
	"wind": {
		"speed": 2.12,
		"deg": 70,
		"gust": 1.89
	},
	"clouds": {
		"all": 0
	},
	"dt": 1743869252,
	"sys": {
		"type": 2,
		"id": 2007821,
		"country": "FR",
		"sunrise": 1743830039,
		"sunset": 1743876859
	},
	"timezone": 7200,
	"id": 2996943,
	"name": "Arrondissement de Lyon",
	"cod": 200
}`

const respBodyMoscou = `{
	"coord": {
		"lon": 37.6156,
		"lat": 55.7522
	},
	"weather": [
		{
		"id": 600,
		"main": "Snow",
		"description": "light snow",
		"icon": "13d"
		}
	],
	"base": "stations",
	"main": {
		"temp": 0.29,
		"feels_like": -4.53,
		"temp_min": -0.76,
		"temp_max": 0.77,
		"pressure": 1007,
		"humidity": 67,
		"sea_level": 1007,
		"grnd_level": 988
	},
	"visibility": 209,
	"wind": {
		"speed": 4.96,
		"deg": 20,
		"gust": 8.88
	},
	"snow": {
		"1h": 0.42
	},
	"clouds": {
		"all": 100
	},
	"dt": 1744101768,
	"sys": {
		"type": 2,
		"id": 2094500,
		"country": "RU",
		"sunrise": 1744080128,
		"sunset": 1744129213
	},
	"timezone": 10800,
	"id": 524901,
	"name": "Moscow",
	"cod": 200
}`

func BootstrapUpstreamServer() *ghttp.Server {
	server := ghttp.NewServer()

	genericHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var response openweather_client.WeatherResponse
		if req.RequestURI == "/weather?appid=weather_api_key&q=Lyon&units=metric" {
			e := json.Unmarshal([]byte(respBodyLyon), &response)
			Expect(e).To(BeNil())
		} else if req.RequestURI == "/weather?appid=weather_api_key&q=Moscou&units=metric" {
			e := json.Unmarshal([]byte(respBodyMoscou), &response)
			Expect(e).To(BeNil())
		}

		respBody, _ := json.Marshal(&response)
		_, e := w.Write(respBody)
		Expect(e).To(BeNil())
	}

	server.RouteToHandler(http.MethodGet, regexp.MustCompile("/weather($|/.*)"), genericHandler)

	return server
}
