# uweather

A CLI tool for querying weather data for cities around the world. Built with Go and uses the Open-Meteo API (free, no API key required).

## Features

- **City Management**: Save cities with custom labels
- **Default Location**: Set a default city for quick weather checks
- **Weather Forecast**: Get current weather or multi-day forecasts (up to 7 days)
- **ASCII Art Display**: Beautiful text-based weather visualization (no emojis)
- **Global Coverage**: Query weather for any city worldwide

## Installation

```bash
# Clone the repository
git clone https://github.com/ugur-claw/uweather.git
cd uweather

# Build the binary
go build -o uweather .

# Optional: Add to PATH
sudo mv uweather /usr/local/bin/
```

## Usage

### Show weather for default location

```bash
uweather
```

### Show weather for a saved location

```bash
uweather home
uweather work
uweather vacation
```

### Show weather for any city (not saved)

```bash
uweather Tokyo
uweather London
uweather "New York"
```

### Add a new city

```bash
uweather add Istanbul --label work
uweather add Amsterdam --label vacation
```

### List saved locations

```bash
uweather locations
```

### Set default location

```bash
uweather default work
```

### Remove a saved location

```bash
uweather remove vacation
```

### Show weather forecast

```bash
# 3-day forecast for default location
uweather --days 3

# 7-day forecast for a saved location
uweather home --days 7
```

## Options

- `--days N` - Number of forecast days (1-7, default: 1)
- `--label name` - Label for a new location (used with `add` command)

## Data Storage

Locations are saved in `~/.uweather/locations.json`:

```json
{
  "locations": [
    {
      "label": "home",
      "city": "Istanbul",
      "lat": 41.0138,
      "lon": 28.9497,
      "country": "Türkiye"
    }
  ],
  "default": "home"
}
```

## Examples

```
$ uweather
┌───────────────────────────────────┐
│        ISTANBUL, TÜRKIYE         │
│                                   │
│             \   |   /             │
│              .--.--.              │
│             /   |   \             │
│                                   │
│             Clear sky             │
│                                   │
│       Temperature: 11.8°C         │
│         Wind: 4.3 km/h SE         │
│           Humidity: 70%           │
│                                   │
└───────────────────────────────────┘
```

```
$ uweather --days 3
┌───────────────────────────────────┐
│         WEATHER FORECAST          │
│        ISTANBUL, TÜRKIYE         │
├───────────────────────────────────┤
│Today  11°-18°  .--.               │
│  Rain: 0.8 mm                     │
│Tomorrow  11°-16°  /_/_/_/        │
│  Rain: 2.7 mm                     │
│Tue  10°-14°  /_/_/_/             │
│  Rain: 10.9 mm                    │
└───────────────────────────────────┘
```

## API

Uses [Open-Meteo API](https://open-meteo.com/) - Free weather API with no API key required.

## License

MIT
