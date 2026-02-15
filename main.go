package main

import (
	"fmt"
	"os"

	"github.com/ugur-claw/uweather/cmd"
	"github.com/ugur-claw/uweather/storage"
)

func main() {
	// Set up custom flag parsing
	// We'll handle flags manually after detecting subcommand

	// Ensure config directory exists
	if err := cmd.EnsureConfigDirCommand(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get all args
	args := os.Args[1:]
	if len(args) == 0 {
		// No args - show weather for default location
		if err := cmd.WeatherCommand("", 1); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Parse common flags first
	daysFlag := 1
	labelFlag := ""
	
	// Look for common flags in the args
	filteredArgs := []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--days" || arg == "-days" {
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &daysFlag)
				i++
			}
		} else if arg == "--label" || arg == "-label" {
			if i+1 < len(args) {
				labelFlag = args[i+1]
				i++
			}
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}

	args = filteredArgs

	// Determine the command based on first positional arg
	if len(args) == 0 {
		// Only flags - show weather for default location
		if err := cmd.WeatherCommand("", daysFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// First arg is the command
	command := args[0]

	switch command {
	case "add":
		// uweather add [city] --label [label]
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: city name required. Usage: uweather add [city] --label [label]\n")
			os.Exit(1)
		}
		city := args[1]
		if labelFlag == "" {
			fmt.Fprintf(os.Stderr, "Error: --label is required when adding a city\n")
			os.Exit(1)
		}
		if err := cmd.AddCommand(city, labelFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return

	case "remove":
		// uweather remove [label]
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: label required. Usage: uweather remove [label]\n")
			os.Exit(1)
		}
		label := args[1]
		if err := cmd.RemoveCommand(label); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return

	case "locations", "list", "ls":
		// uweather locations
		if err := cmd.ListCommand(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return

	case "default":
		// uweather default [label]
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: label required. Usage: uweather default [label]\n")
			os.Exit(1)
		}
		label := args[1]
		if err := cmd.DefaultCommand(label); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return

	case "help", "--help", "-h":
		printHelp()
		return

	default:
		// Treat as city name or label
		arg := args[0]

		// Check if it's a saved label
		locations, _, err := storage.ListLocations()
		if err == nil {
			isLabel := false
			for _, loc := range locations {
				if loc.Label == arg {
					isLabel = true
					break
				}
			}

			if isLabel {
				// It's a saved label - show weather for that location
				if err := cmd.WeatherCommand(arg, daysFlag); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				}
				return
			}
		}

		// Not a saved label - try as city name
		if err := cmd.WeatherByCityCommand(arg, daysFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}
}

func printHelp() {
	fmt.Println(`uweather - Weather CLI Tool

Usage:
  uweather                          Show weather for default location
  uweather [label]                  Show weather for saved location
  uweather [city]                   Show weather for a city
  uweather add [city] --label [name]  Add a new location
  uweather remove [label]           Remove a saved location
  uweather locations                List all saved locations
  uweather default [label]          Set default location

Options:
  --days N     Show N-day forecast (1-7, default: 1)
  --label name Label for a new location

Examples:
  uweather                          # Show weather for default
  uweather home                     # Show weather for 'home'
  uweather Istanbul                 # Show weather for Istanbul
  uweather add Istanbul --label work
  uweather --days 3                 # 3-day forecast for default
`)
}
