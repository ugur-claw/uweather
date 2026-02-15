package models

// Location represents a saved city with label
type Location struct {
	Label   string  `json:"label"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

// LocationsData represents the JSON structure stored in file
type LocationsData struct {
	Locations []Location `json:"locations"`
	Default   string     `json:"default"`
}

// GeocodingResponse represents Open-Meteo Geocoding API response
type GeocodingResponse struct {
	Results []GeocodingResult `json:"results"`
}

type GeocodingResult struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
	Admin1    string  `json:"admin1"` // State/Province
}

// WeatherResponse represents Open-Meteo Weather API response
type WeatherResponse struct {
	CurrentWeather CurrentWeather     `json:"current_weather"`
	Hourly         HourlyWeather      `json:"hourly"`
	Daily          DailyWeather       `json:"daily"`
}

type CurrentWeather struct {
	Temperature   float64 `json:"temperature"`
	Windspeed     float64 `json:"windspeed"`
	Winddirection float64 `json:"winddirection"`
	Weathercode   int     `json:"weathercode"`
	Time          string  `json:"time"`
}

type HourlyWeather struct {
	Time               []string  `json:"time"`
	Temperature_2m     []float64 `json:"temperature_2m"`
	Relativehumidity_2m []int    `json:"relativehumidity_2m"`
}

type DailyWeather struct {
	Time            []string  `json:"time"`
	TemperatureMax  []float64 `json:"temperature_2m_max"`
	TemperatureMin  []float64 `json:"temperature_2m_min"`
	Weathercode     []int     `json:"weathercode"`
	PrecipitationSum []float64 `json:"precipitation_sum"`
}
