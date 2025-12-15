# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go web application that provides a web-based configuration editor using the MVC (Model-View-Controller) pattern. The application serves an HTML form that allows users to modify JSON configuration settings through a web interface, supporting various field types (text, number, boolean) with automatic type conversion and persistence.

## Architecture

The application follows MVC (Model-View-Controller) pattern with clear separation of concerns:

### MVC Layers

**Models (`models/`)**: Data structures and business logic
- `Config` and `ConfigField` structs for ordered field representation
- `ConfigManager` for file operations and data persistence
- Type conversion and form data processing logic
- Automatic migration from legacy to new ordered format
- Field ordering preservation using array-based storage

**Views (`views/`)**: Presentation layer
- `View` struct for template management using external HTML files
- `RenderHome()` method for HTML rendering with `ExecuteTemplate()`
- External HTML template in `templates/index.html`
- Separation of presentation logic from application code

**Controllers (`controllers/`)**: Request handling and orchestration
- `Controller` struct coordinating between model and view
- HTTP request handlers (`HandleHome`, `HandleSave`, `HandleFieldValue`)
- Request validation and error handling with proper HTTP status codes
- Special Boolean API handling: returns HTTP 203 for `false` values, HTTP 200 for `true` values

**Main (`main.go`)**: Application entry point
- Dependency injection and component initialization
- Route registration
- Server startup

### Configuration System

- **Configuration storage**: Uses `options.json` as the configuration file
- **Type system**: Automatically detects and preserves data types (string, int, float, bool)
- **No external dependencies**: Uses only Go standard library

## Development Commands

### Running the Application
```bash
go run main.go
```
The server starts on `http://localhost:8080`

### Building
```bash
go build -o options-api .
```

### Docker Build and Deploy
```bash
# Show all Docker commands
make help

# Build Docker image
make build

# Build and push to Docker Hub
make push

# Run container locally
make run

# Development commands
make dev-run          # Run with volume mount for development
make logs             # View container logs
make shell            # Get interactive shell in container
make stop             # Stop and remove container
make clean            # Clean up images and containers
```

### Module Management
```bash
go mod tidy          # Clean up dependencies
go mod download      # Download dependencies
```

### Testing the Application
```bash
# Test locally
go run main.go

# Test with Docker
make build && make run

# Test with Docker Compose
docker-compose up --build -d

# Access web interface at http://localhost:8080
```

## Environment Variables

The application supports the following environment variables:

- **`CONFIG_FILE`**: Path to the configuration file (default: `options.json`)
- **Usage**:
  ```bash
  CONFIG_FILE=/custom/path/config.json go run main.go
  ```

## Docker Compose Support

For easier local development and deployment, the project includes a `docker-compose.yml` file:

### Quick Start with Docker Compose
```bash
# Build and run all services
docker-compose up --build -d

# View logs
docker-compose logs -f

# Stop and remove containers
docker-compose down
```

### Docker Compose Features
- **Volume Mounting**: Persistent configuration hot-reloading
- **Health Checks**: Automatic health monitoring with wget-based checks
- **Auto-restart**: Container restarts on failure (`unless-stopped`)
- **Production Ready**: Optimized for deployment with proper health checks

## Configuration File Format

The `options.json` file uses an ordered array structure to maintain consistent field ordering:

**Current Format (New):**
```json
{
  "fields": [
    {"name": "fieldText", "value": "Текстовое значение"},
    {"name": "intData", "value": 100500},
    {"name": "boolValue", "value": true},
    {"name": "isDisabled", "value": false},
    {"name": "eventDate", "value": "2025-12-25 15:30:00"}
  ]
}
```

**Legacy Format (Still Supported):**
```json
{
  "fieldText": "Текстовое значение",
  "intData": 100500,
  "boolValue": true
}
```

The application automatically migrates legacy configurations to the new ordered format while preserving field order using predefined sequence: `fieldText`, `intData`, `boolValue`, with additional fields appended in their original order.

### Supported Field Types
- **Text fields**: Rendered as HTML text inputs, support date format detection for calendar integration
- **Dates**: Special processing for `YYYY-MM-DD HH:MM:SS` format with API status code differentiation
- **Numbers**: Rendered as HTML number inputs (supports int, int64, float64)
- **Booleans**: Rendered as HTML checkboxes with special API status codes
- **Complex types**: Displayed as read-only values

## Key Features

- **Dynamic form generation**: Automatically generates form fields based on JSON structure
- **Type preservation**: Maintains original data types when saving configuration
- **Field ordering**: Consistent field display order preserved across page reloads and form submissions
- **Responsive design**: Mobile-friendly CSS with modern styling
- **Form validation**: Client-side validation for unsaved changes with beforeunload confirmation
- **Success feedback**: Visual confirmation when settings are saved
- **Default configuration**: Auto-creates config file with sample values if missing
- **MVC architecture**: Clean separation of concerns for maintainability
- **Template separation**: External HTML templates for easier frontend development
- **Backward compatibility**: Automatic migration from legacy configuration formats
- **Docker support**: Multi-stage Docker builds with optimized image size
- **Containerization**: Ready for deployment with Docker and Kubernetes support
- **Calendar functionality**: Interactive date/time picker with "Сейчас" button for text fields containing date formats (YYYY-MM-DD HH:MM:SS)
- **Enhanced API responses**: Special HTTP status codes for Boolean fields (203 for false, 200 for true) and date fields (203 for past dates, 200 for future/present dates)

## Project Structure

```
.
├── main.go              # Application entry point and dependency injection
├── go.mod              # Go module definition
├── Dockerfile           # Multi-stage Docker build configuration
├── Makefile            # Build and deployment automation
├── .dockerignore       # Files excluded from Docker build
├── options.json         # Configuration file with ordered fields (auto-created)
├── models/              # Model layer - data structures and business logic
│   └── config.go        # Config management, ordered fields, and format migration
├── views/               # View layer - presentation logic
│   └── template.go      # HTML template rendering with external files
├── templates/           # External HTML templates
│   └── index.html       # Main configuration form template
├── controllers/         # Controller layer - request handling
│   └── controller.go    # HTTP handlers and orchestration
├── CLAUDE.md           # This file
└── DOCKER.md           # Docker deployment guide
```

## Development Notes

- **Module name**: `options-api`
- **Go version**: 1.21+
- **Dependencies**: Only Go standard library
- **Architecture**: Clean MVC pattern with dependency injection
- **Configuration**: Stored in ordered JSON array format with backward compatibility
- **Error handling**: Comprehensive error handling with proper HTTP status codes
- **Language**: All user-facing text is in Russian
- **Port**: Server runs on port 8080 by default
- **Template system**: External HTML templates with `ExecuteTemplate()` method
- **Field ordering**: Guaranteed consistent field order using array-based configuration storage

## MVC Implementation Details

### Dependency Flow
```go
// In main.go - Dependency injection
configManager := models.NewConfigManager(configFile)
view, _ := views.NewView()
controller := controllers.NewController(configManager, view)
```

### Request Flow
1. HTTP request → Controller handler
2. Controller → Model (data operations)
3. Model → Controller (data/results)
4. Controller → View (template rendering)
5. View → HTTP response

### Data Flow
- **Form submission**: `r.Form` → `Controller` → `Model.UpdateConfigFromForm()` → `Model.SaveConfig()`
- **Page rendering**: `Model.LoadConfig()` → `Model.GetFields()` → `Controller` → `View.RenderHome()`

### API Endpoints

The application provides several REST API endpoints:

**Web Interface:**
- `GET /` - Main configuration form
- `POST /save` - Save configuration changes

**Field Value API:**
- `GET /{fieldName}` - Retrieve individual field values
  - Returns `text/plain` content type
  - **Boolean fields**: HTTP 200 for `true`, HTTP 203 for `false`
  - **Date fields**: HTTP 200 for future/present dates, HTTP 203 for past dates (format: `YYYY-MM-DD HH:MM:SS`)
  - **Other fields**: HTTP 200 with field value
  - **Non-existent fields**: HTTP 404 with error message

**Enhanced Date/Time Processing:**
- Date fields in format `YYYY-MM-DD HH:MM:SS` receive special status code handling
- **Past dates**: Return HTTP 203 Non-Authoritative Information
- **Future/Present dates**: Return HTTP 200 OK
- This provides semantic meaning about date expiration status

## Frontend Features

### Calendar Integration
The application includes an intelligent date/time picker system:

**Automatic Detection:**
- Scans text fields for date format `YYYY-MM-DD HH:MM:SS`
- Dynamically shows/hides calendar button (📅) based on content
- Preserves existing field values and formats

**Interactive Calendar:**
- Modal dialog with date and time inputs
- **"Сейчас" button**: Sets current date/time instantly
- Manual date/time selection with precision to seconds
- Responsive design with modern styling

**Usage:**
1. Text field contains date in correct format → calendar button appears
2. Click calendar button → modal opens with current field values
3. Use "Сейчас" for instant current time or select manually
4. Apply updates field with formatted date/time
5. Form change detection tracks modifications for unsaved changes warnings

### Field Ordering Architecture

The application ensures consistent field ordering through several key architectural decisions:

**Data Structure Design:**
- Uses `[]ConfigField` instead of `map[string]interface{}` to preserve order
- Each field stores `name` and `value` in a structured array
- Order is maintained from configuration file through to HTML rendering

**Migration Strategy:**
- Legacy configurations are automatically converted to ordered format
- Predefined field order sequence: `fieldText`, `intData`, `boolValue`
- Additional fields (`isDisabled`, `eventDate`, etc.) are appended in the order they appear in the original file

**Template Integration:**
- External HTML templates loaded via `template.ParseFiles()`
- Field iteration uses Go template `range` directive over ordered array
- Consistent rendering guaranteed by array-based iteration