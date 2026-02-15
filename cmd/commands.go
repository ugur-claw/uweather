package cmd

import (
	"fmt"

	"github.com/ugur-claw/uweather/api"
	"github.com/ugur-claw/uweather/models"
	"github.com/ugur-claw/uweather/storage"
	"github.com/ugur-claw/uweather/ui"
)

// AddCommand adds a new city location
func AddCommand(city, label string) error {
	if city == "" {
		return fmt.Errorf("city name is required")
	}
	if label == "" {
		return fmt.Errorf("label is required")
	}

	client := api.NewClient()
	result, err := client.Geocoding(city)
	if err != nil {
		return err
	}

	err = storage.AddLocation(label, result.Name, result.Latitude, result.Longitude, result.Country)
	if err != nil {
		return err
	}

	fmt.Printf("Added: %s (%s, %.4f, %.4f) with label '%s'\n",
		result.Name, result.Country, result.Latitude, result.Longitude, label)

	return nil
}

// RemoveCommand removes a city by label
func RemoveCommand(label string) error {
	if label == "" {
		return fmt.Errorf("label is required")
	}

	err := storage.RemoveLocation(label)
	if err != nil {
		return err
	}

	fmt.Printf("Removed: %s\n", label)
	return nil
}

// ListCommand lists all saved locations
func ListCommand() error {
	locations, defaultLabel, err := storage.ListLocations()
	if err != nil {
		return err
	}

	if len(locations) == 0 {
		fmt.Println("No locations saved. Add a location with 'uweather add [city] --label [label]'")
		return nil
	}

	fmt.Println("Saved locations:")
	fmt.Println("---------------")
	for _, loc := range locations {
		marker := " "
		if loc.Label == defaultLabel {
			marker = "*"
		}
		fmt.Printf("%s %s -> %s, %s\n", marker, loc.Label, loc.City, loc.Country)
	}
	fmt.Println("\n* = default location")

	return nil
}

// DefaultCommand sets the default location
func DefaultCommand(label string) error {
	if label == "" {
		return fmt.Errorf("label is required")
	}

	err := storage.SetDefaultLocation(label)
	if err != nil {
		return err
	}

	fmt.Printf("Default location set to: %s\n", label)
	return nil
}

// WeatherCommand fetches and displays weather for a label or default
func WeatherCommand(label string, days int) error {
	var location *models.Location
	var err error

	if label == "" {
		// Use default location
		location, err = storage.GetDefaultLocation()
		if err != nil {
			return err
		}
	} else {
		// Use specified label
		location, err = storage.GetLocation(label)
		if err != nil {
			return err
		}
	}

	client := api.NewClient()
	weather, err := client.GetWeather(location.Lat, location.Lon, days)
	if err != nil {
		return err
	}

	// Display weather
	ui.DisplayWeather(location, weather, days)

	return nil
}

// WeatherByCityCommand fetches weather for a city without saving
func WeatherByCityCommand(city string, days int) error {
	if city == "" {
		return fmt.Errorf("city name is required")
	}

	client := api.NewClient()
	result, err := client.Geocoding(city)
	if err != nil {
		return err
	}

	weather, err := client.GetWeather(result.Latitude, result.Longitude, days)
	if err != nil {
		return err
	}

	location := &models.Location{
		City:    result.Name,
		Country: result.Country,
		Lat:     result.Latitude,
		Lon:     result.Longitude,
	}

	ui.DisplayWeather(location, weather, days)

	return nil
}

// EnsureConfigDir ensures the config directory exists
func EnsureConfigDirCommand() error {
	return storage.EnsureConfigDir()
}

// GetDefaultLabel returns the default label
func GetDefaultLabel() (string, error) {
	_, defaultLabel, err := storage.ListLocations()
	return defaultLabel, err
}
