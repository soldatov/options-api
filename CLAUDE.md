# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a simple Go web application that provides a web-based configuration editor. The application serves an HTML form that allows users to modify JSON configuration settings through a web interface, supporting various field types (text, number, boolean) with automatic type conversion and persistence.

## Architecture

- **Single-file application**: All logic is contained in `main.go`
- **Configuration storage**: Uses `options.json` as the configuration file
- **Web interface**: Self-contained HTML/CSS/JS template embedded in Go code
- **Type system**: Automatically detects and preserves data types (string, int, float, bool)
- **No external dependencies**: Uses only Go standard library

### Core Components

- `main()`: Initializes server, loads config, sets up HTTP routes
- `homeHandler()`: Renders the configuration form (GET /)
- `saveHandler()`: Processes form submissions and updates config (POST /save)
- `readConfig()`/`saveConfig()`: JSON file operations with type preservation
- `loadConfig()`: Creates default config if file doesn't exist
- `htmlTemplate()`: Embedded HTML template with styling and JavaScript

## Development Commands

### Running the Application
```bash
go run main.go
```
The server starts on `http://localhost:8080`

### Building
```bash
go build -o options-api main.go
```

### Testing the Application
```bash
# Test the web interface manually by accessing http://localhost:8080
# No automated tests are currently implemented
```

## Configuration File Format

The `options.json` file stores configuration as a flat JSON object:
```json
{
  "fieldText": "Текстовое значение",
  "intData": 100500,
  "boolValue": true
}
```

### Supported Field Types
- **Text fields**: Rendered as HTML text inputs
- **Numbers**: Rendered as HTML number inputs (supports int, int64, float64)
- **Booleans**: Rendered as HTML checkboxes
- **Complex types**: Displayed as read-only values

## Key Features

- **Dynamic form generation**: Automatically generates form fields based on JSON structure
- **Type preservation**: Maintains original data types when saving configuration
- **Responsive design**: Mobile-friendly CSS with modern styling
- **Form validation**: Client-side validation for unsaved changes
- **Success feedback**: Visual confirmation when settings are saved
- **Default configuration**: Auto-creates config file with sample values if missing

## Project Structure

```
.
├── main.go          # Complete application source code
├── options.json     # Configuration file (auto-created)
└── CLAUDE.md        # This file
```

## Development Notes

- The application uses only the Go standard library - no external dependencies
- HTML templates are embedded as strings in the Go source code
- Configuration is stored in a flat JSON structure (no nested objects)
- The application includes comprehensive error handling for file operations
- All text is in Russian language as per the original implementation
- Server runs on port 8080 by default