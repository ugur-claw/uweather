package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ugur-claw/uweather/models"
)

// Client handles Open-Meteo API requests
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

// Geocoding searches for a city and returns coordinates
func (c *Client) Geocoding(query string) (*models.GeocodingResult, error) {
	results, err := c.GeocodingMulti(query)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("city not found: %s", query)
	}
	// If only one result, return it directly
	if len(results) == 1 {
		return &results[0], nil
	}
	// Multiple results - return the most likely one (first one)
	// The caller should use GeocodingMulti to get all options
	return &results[0], nil
}

// GeocodingMulti searches for a city and returns all matching coordinates
func (c *Client) GeocodingMulti(query string) ([]models.GeocodingResult, error) {
	// Encode the query
	encodedQuery := url.QueryEscape(query)
	// Request more results to allow selection
	url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=10&language=en&format=json", encodedQuery)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("geocoding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var geocodingResp models.GeocodingResponse
	if err := json.Unmarshal(body, &geocodingResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(geocodingResp.Results) == 0 {
		return nil, fmt.Errorf("city not found: %s", query)
	}

	return geocodingResp.Results, nil
}

// GetGeocodingResults returns all geocoding results for a query
func (c *Client) GetGeocodingResults(query string) ([]models.GeocodingResult, error) {
	return c.GeocodingMulti(query)
}

// GetWeather fetches weather data for given coordinates
func (c *Client) GetWeather(lat, lon float64, days int) (*models.WeatherResponse, error) {
	// Limit days to max 7
	if days > 7 {
		days = 7
	}
	if days < 1 {
		days = 1
	}

	// Build hourly params for today
	hourlyParams := "temperature_2m,relativehumidity_2m"
	dailyParams := "temperature_2m_max,temperature_2m_min,weathercode,precipitation_sum"

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current_weather=true&hourly=%s&daily=%s&timezone=auto&forecast_days=%d",
		lat, lon, hourlyParams, dailyParams, days)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("weather request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var weatherResp models.WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &weatherResp, nil
}

// GetWeatherCodeDescription returns description for WMO weather code
func GetWeatherCodeDescription(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1:
		return "Mainly clear"
	case 2:
		return "Partly cloudy"
	case 3:
		return "Overcast"
	case 45, 48:
		return "Foggy"
	case 51, 53, 55:
		return "Drizzle"
	case 56, 57:
		return "Freezing drizzle"
	case 61, 63, 65:
		return "Rain"
	case 66, 67:
		return "Freezing rain"
	case 71, 73, 75:
		return "Snow"
	case 77:
		return "Snow grains"
	case 80, 81, 82:
		return "Rain showers"
	case 85, 86:
		return "Snow showers"
	case 95:
		return "Thunderstorm"
	case 96, 99:
		return "Thunderstorm with hail"
	default:
		return "Unknown"
	}
}

// GetWeatherArt returns ASCII art based on weather code
func GetWeatherArt(code int) string {
	switch code {
	case 0: // Clear sky
		return `      \   |   /    
       .-.       
    - (   ) -   
       -'-      
      /   \     `
	case 1: // Mainly clear
		return `      \   |   /    
   .--.  .-     
  (    ).(___)  
  (___.__)__ )  `
	case 2: // Partly cloudy
		return `      \  |  /     
  .--. .--.     
 .(    ).(___. 
(___.__)____)  `
	case 3: // Overcast
		return `             
  .--.       
 .(   ).     
(___.__)     `
	case 45, 48: // Fog
		return `      _ - _ -   
     _ - _ - _  
    _ - _ - _   
     _ - _ -    `
	case 51, 53, 55: // Drizzle
		return `       .--.    
    .-(    ).   
   (___.__)__)  
    ' ' ' ' '    `
	case 56, 57: // Freezing drizzle
		return `       .--.    
    .-(    ).   
   (___.__)__)  
    ', ', ', '   `
	case 61, 63, 65: // Rain
		return `       .--.    
    .-(    ).   
   (___.__)__)  
   ' ' ' ' '    `
	case 66, 67: // Freezing rain
		return `       .--.    
    .-(    ).   
   (___.__)__)  
   ,,,,,,,,,    `
	case 71, 73, 75: // Snow
		return `       .--.    
    .-(    ).   
   (___.__)__)  
    * * * *     `
	case 77: // Snow grains
		return `       .--.    
    .-(    ).   
   (___.__)__)  
     . . . .    `
	case 80, 81, 82: // Rain showers
		return `     _` + "`" + `QQQ     
       .--.    
    .-(    ).   
   (___.__)__)  
   ' ' ' ' '    `
	case 85, 86: // Snow showers
		return `     ***     
       .--.    
    .-(    ).   
   (___.__)__)  
    *  *  *     `
	case 95: // Thunderstorm
		return `     /` + `_` + `/` + `_` + `/` + `_` + `/  
       .--.    
    .-(    ).   
   (___.__)__)  
     /` + `_` + `/` + `_` + `/    `
	case 96, 99: // Thunderstorm with hail
		return `     /` + `_` + `/` + `_` + `/` + `_` + `/  
       .--.    
    .-(    ).   
   (___.__)__)  
   ,***,***,    `
	default:
		return `      \   |   /    
       .-.       
    - (   ) -   
       -'-      
      /   \     `
	}
}

// FormatWindDirection converts wind direction degrees to cardinal direction
func FormatWindDirection(degrees float64) string {
	dirs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int((degrees + 22.5) / 45) % 8
	return dirs[index]
}

// FormatCityName formats city name for display
func FormatCityName(city, country, admin1 string) string {
	parts := []string{strings.ToUpper(city)}
	if country != "" {
		parts = append(parts, strings.ToUpper(country))
	} else if admin1 != "" {
		parts = append(parts, strings.ToUpper(admin1))
	}
	return strings.Join(parts, ", ")
}

// GetWeatherEmoji returns emoji for weather code (for ASCII tables)
func GetWeatherEmoji(code int) string {
	switch code {
	case 0:
		return "☀️"
	case 1:
		return "🌤️"
	case 2:
		return "⛅"
	case 3:
		return "☁️"
	case 45, 48:
		return "🌫️"
	case 51, 53, 55:
		return "🌧️"
	case 56, 57:
		return "🌨️"
	case 61, 63, 65:
		return "🌧️"
	case 66, 67:
		return "🌨️"
	case 71, 73, 75:
		return "❄️"
	case 77:
		return "🌨️"
	case 80, 81, 82:
		return "🌦️"
	case 85, 86:
		return "🌨️"
	case 95:
		return "⛈️"
	case 96, 99:
		return "⛈️"
	default:
		return "🌡️"
	}
}
