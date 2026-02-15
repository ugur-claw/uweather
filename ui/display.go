package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ugur-claw/uweather/api"
	"github.com/ugur-claw/uweather/models"
)

// DisplayWeather displays weather information with ASCII art
func DisplayWeather(location *models.Location, weather *models.WeatherResponse, days int) {
	cityName := api.FormatCityName(location.City, location.Country, "")

	if days == 1 {
		// Single day display (current weather)
		displayCurrentWeather(cityName, weather)
	} else {
		// Multi-day forecast - use ASCII table
		displayForecastTable(cityName, weather, days)
	}
}

func displayCurrentWeather(cityName string, weather *models.WeatherResponse) {
	current := weather.CurrentWeather
	art := api.GetWeatherArt(current.Weathercode)
	desc := api.GetWeatherCodeDescription(current.Weathercode)
	windDir := api.FormatWindDirection(current.Winddirection)

	// Get humidity for current hour
	humidity := 0
	if len(weather.Hourly.Relativehumidity_2m) > 0 {
		currentHour := time.Now().Hour()
		if currentHour < len(weather.Hourly.Relativehumidity_2m) {
			humidity = weather.Hourly.Relativehumidity_2m[currentHour]
		}
	}

	// Calculate box width
	width := 37

	fmt.Println("┌" + strings.Repeat("─", width-2) + "┐")
	fmt.Printf("│%s%s%s│\n", strings.Repeat(" ", (width-2-len(cityName))/2), cityName, strings.Repeat(" ", (width-2-len(cityName)+1)/2))
	fmt.Println("│" + strings.Repeat(" ", width-2) + "│")

	// Split and display art
	artLines := strings.Split(art, "\n")
	for _, line := range artLines {
		if strings.TrimSpace(line) != "" {
			padding := (width - 2 - len(strings.TrimSpace(line))) / 2
			fmt.Printf("│%s%s%s│\n", strings.Repeat(" ", padding), strings.TrimSpace(line), strings.Repeat(" ", width-2-len(strings.TrimSpace(line))-padding))
		}
	}

	fmt.Println("│" + strings.Repeat(" ", width-2) + "│")

	// Weather description
	descPadding := (width - 2 - len(desc)) / 2
	fmt.Printf("│%s%s%s│\n", strings.Repeat(" ", descPadding), desc, strings.Repeat(" ", width-2-len(desc)-descPadding))

	fmt.Println("│" + strings.Repeat(" ", width-2) + "│")

	// Temperature
	tempLine := fmt.Sprintf("Temperature: %.1f°C", current.Temperature)
	tempPadding := (width - 2 - len(tempLine)) / 2
	fmt.Printf("│%s%s%s│\n", strings.Repeat(" ", tempPadding), tempLine, strings.Repeat(" ", width-2-len(tempLine)-tempPadding))

	// Wind
	windLine := fmt.Sprintf("Wind: %.1f km/h %s", current.Windspeed, windDir)
	windPadding := (width - 2 - len(windLine)) / 2
	fmt.Printf("│%s%s%s│\n", strings.Repeat(" ", windPadding), windLine, strings.Repeat(" ", width-2-len(windLine)-windPadding))

	// Humidity
	if humidity > 0 {
		humLine := fmt.Sprintf("Humidity: %d%%", humidity)
		humPadding := (width - 2 - len(humLine)) / 2
		fmt.Printf("│%s%s%s│\n", strings.Repeat(" ", humPadding), humLine, strings.Repeat(" ", width-2-len(humLine)-humPadding))
	}

	fmt.Println("│" + strings.Repeat(" ", width-2) + "│")
	fmt.Println("└" + strings.Repeat("─", width-2) + "┘")
}

func displayForecast(cityName string, weather *models.WeatherResponse, days int) {
	// Header
	width := 37
	fmt.Println("┌" + strings.Repeat("─", width-2) + "┐")
	title := "WEATHER FORECAST"
	titlePadding := (width - 2 - len(title)) / 2
	fmt.Printf("│%s%s%s│\n", strings.Repeat(" ", titlePadding), title, strings.Repeat(" ", width-2-len(title)-titlePadding))

	cityPadding := (width - 2 - len(cityName)) / 2
	fmt.Printf("│%s%s%s│\n", strings.Repeat(" ", cityPadding), cityName, strings.Repeat(" ", width-2-len(cityName)-cityPadding))
	fmt.Println("├" + strings.Repeat("─", width-2) + "┤")

	// Display each day
	daily := weather.Daily
	for i := 0; i < days && i < len(daily.Time); i++ {
		// Parse date
		date, err := time.Parse("2006-01-02", daily.Time[i])
		if err != nil {
			continue
		}

		dayName := date.Format("Mon")
		if i == 0 {
			dayName = "Today"
		} else if i == 1 {
			dayName = "Tomorrow"
		}

		tempMax := daily.TemperatureMax[i]
		tempMin := daily.TemperatureMin[i]
		code := daily.Weathercode[i]
		precip := daily.PrecipitationSum[i]

		art := api.GetWeatherArt(code)
		artLines := strings.Split(art, "\n")
		artLine := ""
		for _, line := range artLines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				artLine = trimmed
				break
			}
		}

		// Weather line
		weatherLine := fmt.Sprintf("%s  %.0f°-%.0f°  %s", dayName, tempMin, tempMax, artLine)
		fmt.Printf("│%s%s│\n", weatherLine, strings.Repeat(" ", width-2-len(weatherLine)))

		// Precipitation
		if precip > 0 {
			precipLine := fmt.Sprintf("  Rain: %.1f mm", precip)
			fmt.Printf("│%s%s│\n", precipLine, strings.Repeat(" ", width-2-len(precipLine)))
		}
	}

	fmt.Println("└" + strings.Repeat("─", width-2) + "┘")
}

func displayForecastTable(cityName string, weather *models.WeatherResponse, days int) {
	// Header with city name
	fmt.Println()
	fmt.Println("┌" + strings.Repeat("─", 43) + "┐")
	title := "WEATHER FORECAST - " + cityName
	titlePadding := (44 - len(title)) / 2
	fmt.Printf("│%s%s%s│\n", strings.Repeat(" ", titlePadding), title, strings.Repeat(" ", 44-len(title)-titlePadding))
	fmt.Println("├" + strings.Repeat("─", 15) + "┬" + strings.Repeat("─", 11) + "┬" + strings.Repeat("─", 10) + "┬" + strings.Repeat("─", 6) + "┤")
	fmt.Printf("│%s│%s│%s│%s│\n", centerText(" Day ", 15), centerText("  Temp  ", 11), centerText("  Wind  ", 10), centerText(" Status ", 6))
	fmt.Println("├" + strings.Repeat("─", 15) + "┼" + strings.Repeat("─", 11) + "┼" + strings.Repeat("─", 10) + "┼" + strings.Repeat("─", 6) + "┤")

	// Current weather for wind info
	current := weather.CurrentWeather
	windSpeed := current.Windspeed

	// Display each day
	daily := weather.Daily
	for i := 0; i < days && i < len(daily.Time); i++ {
		// Parse date
		date, err := time.Parse("2006-01-02", daily.Time[i])
		if err != nil {
			continue
		}

		dayName := date.Format("Mon Jan 2")
		if i == 0 {
			dayName = "Today"
		} else if i == 1 {
			dayName = "Tomorrow"
		}

		tempMax := daily.TemperatureMax[i]
		tempMin := daily.TemperatureMin[i]
		code := daily.Weathercode[i]

		temp := fmt.Sprintf("%.0f°-%.0f°C", tempMin, tempMax)
		wind := fmt.Sprintf("%.0fkm/h", windSpeed)
		status := api.GetWeatherEmoji(code)

		fmt.Printf("│%s│%s│%s│%s│\n",
			centerText(" "+dayName, 15),
			centerText(temp, 11),
			centerText(wind, 10),
			centerText(" "+status, 6))
	}

	fmt.Println("└" + strings.Repeat("─", 15) + "┴" + strings.Repeat("─", 11) + "┴" + strings.Repeat("─", 10) + "┴" + strings.Repeat("─", 6) + "┘")
	fmt.Println()
}

func centerText(text string, width int) string {
	padding := width - len(text)
	if padding <= 0 {
		return text[:width]
	}
	left := padding / 2
	right := padding - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}
