package entities

type WeatherType string

const (
	Clear        WeatherType = "clear"
	Clouds       WeatherType = "clouds"
	Drizzle      WeatherType = "drizzle"
	Rain         WeatherType = "rain"
	Snow         WeatherType = "snow"
	Thunderstorm WeatherType = "thunderstorm"
)

type WeatherQuery struct {
	Town    string
	TempMin int64
	Type    WeatherType
}
