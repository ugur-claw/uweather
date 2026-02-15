package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ugur-claw/uweather/models"
)

const (
	configDir  = ".uweather"
	configFile = "locations.json"
)

// GetConfigPath returns the path to the config directory
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %w", err)
	}
	return filepath.Join(homeDir, configDir), nil
}

// GetConfigFilePath returns the full path to the config file
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configFile), nil
}

// EnsureConfigDir ensures the config directory exists
func EnsureConfigDir() error {
	configDir, err := GetConfigPath()
	if err != nil {
		return err
	}
	return os.MkdirAll(configDir, 0755)
}

// LoadLocations loads locations from the config file
func LoadLocations() (*models.LocationsData, error) {
	configFile, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return empty data
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return &models.LocationsData{
			Locations: []models.Location{},
			Default:   "",
		}, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var locations models.LocationsData
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &locations, nil
}

// SaveLocations saves locations to the config file
func SaveLocations(data *models.LocationsData) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configFile, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(configFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// AddLocation adds a new location
func AddLocation(label, city string, lat, lon float64, country string) error {
	data, err := LoadLocations()
	if err != nil {
		return err
	}

	// Check if label already exists
	for _, loc := range data.Locations {
		if loc.Label == label {
			return fmt.Errorf("label '%s' already exists. Use a different label or remove it first.", label)
		}
	}

	newLocation := models.Location{
		Label:   label,
		City:    city,
		Lat:     lat,
		Lon:     lon,
		Country: country,
	}

	data.Locations = append(data.Locations, newLocation)

	// If this is the first location, set it as default
	if data.Default == "" {
		data.Default = label
	}

	return SaveLocations(data)
}

// RemoveLocation removes a location by label
func RemoveLocation(label string) error {
	data, err := LoadLocations()
	if err != nil {
		return err
	}

	found := false
	newLocations := make([]models.Location, 0)
	for _, loc := range data.Locations {
		if loc.Label == label {
			found = true
			continue
		}
		newLocations = append(newLocations, loc)
	}

	if !found {
		return fmt.Errorf("label '%s' not found", label)
	}

	data.Locations = newLocations

	// If removed was default, clear default
	if data.Default == label {
		if len(data.Locations) > 0 {
			data.Default = data.Locations[0].Label
		} else {
			data.Default = ""
		}
	}

	return SaveLocations(data)
}

// GetLocation returns a location by label
func GetLocation(label string) (*models.Location, error) {
	data, err := LoadLocations()
	if err != nil {
		return nil, err
	}

	for _, loc := range data.Locations {
		if loc.Label == label {
			return &loc, nil
		}
	}

	return nil, fmt.Errorf("label '%s' not found", label)
}

// GetDefaultLocation returns the default location
func GetDefaultLocation() (*models.Location, error) {
	data, err := LoadLocations()
	if err != nil {
		return nil, err
	}

	if data.Default == "" {
		return nil, fmt.Errorf("no default location set. Use 'uweather default [label]' to set one")
	}

	return GetLocation(data.Default)
}

// SetDefaultLocation sets the default location by label
func SetDefaultLocation(label string) error {
	// First check if label exists
	_, err := GetLocation(label)
	if err != nil {
		return err
	}

	data, err := LoadLocations()
	if err != nil {
		return err
	}

	data.Default = label
	return SaveLocations(data)
}

// ListLocations returns all saved locations
func ListLocations() ([]models.Location, string, error) {
	data, err := LoadLocations()
	if err != nil {
		return nil, "", err
	}

	return data.Locations, data.Default, nil
}
