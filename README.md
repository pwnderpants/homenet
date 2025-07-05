# HTMX + Go + Tailwind CSS Boilerplate

A modern web application boilerplate built with HTMX for dynamic interactions, Go for the backend, and Tailwind CSS for styling with dark mode support.

## Features

- 🚀 **HTMX Integration**: Dynamic interactions without writing JavaScript
- ⚡ **Go Backend**: Fast, simple, and reliable web server
- 🎨 **Tailwind CSS**: Utility-first CSS framework
- 🌙 **Dark Theme**: Fixed dark theme design
- 📱 **Responsive Design**: Mobile-first responsive layout
- 🎬 **Movie Board**: Interactive movie list management with HTMX
- 📝 **Structured Logging**: Comprehensive logging system with configurable levels
- ⚙️ **Configuration System**: JSON-based configuration file with automatic defaults

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go      # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go    # Configuration management
│   ├── handlers/
│   │   ├── handlers.go  # HTTP request handlers
│   │   ├── ollama.go    # Ollama AI integration
│   │   ├── types.go     # Data structures
│   │   └── declarations.go # Constants and configurations
│   ├── database/
│   │   └── database.go  # Database operations
│   ├── logger/
│   │   └── logger.go    # Structured logging system
│   ├── server/
│   │   ├── server.go    # Server configuration
│   │   └── declarations.go # Server types
│   └── templates/
│       └── templates.go # Template management
├── web/
│   ├── templates/
│   │   ├── index.html   # Homepage template
│   │   ├── movie-board.html # Movie management
│   │   ├── tv-shows-board.html # TV show management
│   │   └── ai.html      # AI chat interface
│   └── static/
│       ├── css/
│       │   └── custom.css
│       └── js/
├── go.mod               # Go module file
├── Makefile            # Build and development commands
├── .air.toml           # Hot reload configuration
└── README.md           # This file
```

## Prerequisites

- Go 1.24.3 or later
- A modern web browser

## Getting Started

1. **Clone or download the project**

2. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

3. **Open your browser** and navigate to:
   ```
   http://localhost:8080
   ```

## Configuration

The application uses a JSON configuration file located at `~/.config/homenet/config.json`. This file is automatically created with default values when the application starts for the first time.

### Configuration File Location

```
~/.config/homenet/config.json
```

### Default Configuration

```json
{
  "ollama": {
    "host": "http://chadgpt.gotpwnd.org:11434",
    "model_name": "llama3.2:latest"
  },
  "logging": {
    "level": "INFO"
  },
  "server": {
    "port": "8080"
  }
}
```

### Configuration Options

#### Ollama Settings
- **`ollama.host`**: The Ollama server URL (default: `http://chadgpt.gotpwnd.org:11434`)
- **`ollama.model_name`**: The Ollama model to use for AI queries (default: `llama3.2:latest`)

#### Logging Settings
- **`logging.level`**: Log level for the application
  - `DEBUG`: Detailed debug information
  - `INFO`: General information about application flow
  - `WARN`: Warning messages for potentially harmful situations
  - `ERROR`: Error messages for error conditions

#### Server Settings
- **`server.port`**: Port number for the web server (default: `8080`)

### Modifying Configuration

Simply edit the `~/.config/homenet/config.json` file to change any settings. The application will read the updated configuration on the next startup.

Example configuration changes:

```json
{
  "ollama": {
    "host": "http://localhost:11434",
    "model_name": "llama3.1:latest"
  },
  "logging": {
    "level": "DEBUG"
  },
  "server": {
    "port": "3000"
  }
}
```

## Usage

### Development

For development, you can run the server with hot reloading using tools like `air`:

```bash
# Install air (optional)
go install github.com/cosmtrek/air@latest

# Run with air for hot reloading
air
```

### Production

For production deployment:

```bash
# Build the binary
go build -o homenet cmd/server/main.go

# Run the binary
./homenet
```

The server will use the port specified in your configuration file (`~/.config/homenet/config.json`).

## Features Explained

### HTMX Integration

The application demonstrates HTMX functionality with the Movie Board:

- **Add Movies**: Form submission with HTMX to add movies to the list
- **Delete Movies**: HTMX-powered delete functionality
- **Server-side rendering** with Go templates
- **Static file serving** for CSS and JavaScript
- **Modular handler structure** for easy HTMX endpoint addition

### Dark Theme

- **Fixed Dark Design**: Application uses a consistent dark theme
- **Optimized Colors**: Carefully selected dark colors for optimal readability
- **Custom Scrollbar**: Dark-themed scrollbar that matches the overall design

### Tailwind CSS

- **Utility Classes**: Rapid UI development with utility classes
- **Dark Mode Variants**: Automatic dark mode styling with `dark:` prefix
- **Responsive Design**: Mobile-first responsive breakpoints
- **Custom Styling**: Custom scrollbar styling for dark mode

### Configuration System

The application uses a JSON-based configuration system that automatically creates a default configuration file at `~/.config/homenet/config.json` on first run. This system provides:

- **Automatic Setup**: Creates configuration directory and file if they don't exist
- **Default Values**: Sensible defaults for all configuration options
- **Runtime Configuration**: No need to recompile to change settings
- **Centralized Settings**: All application settings in one place

The configuration system supports:
- **Ollama Integration**: Host URL and model name configuration
- **Logging Levels**: Configurable log verbosity
- **Server Settings**: Port and other server configuration
- **Extensible Design**: Easy to add new configuration options

### Structured Logging

The application includes a comprehensive logging system that tracks:

- **Database Operations**: Movie and TV show additions, updates, and deletions
- **AI Queries**: Ollama API requests and responses
- **Template Loading**: Template compilation and loading
- **Server Events**: Server startup and configuration
- **Error Handling**: Detailed error logging with context

Log messages include timestamps and log levels for easy debugging and monitoring.

#### Backend Logging (Go)

The log level is configured in the `~/.config/homenet/config.json` file under the `logging.level` field. Available levels are:

- `DEBUG`: Detailed debug information
- `INFO`: General information about application flow  
- `WARN`: Warning messages for potentially harmful situations
- `ERROR`: Error messages for error conditions

To change the log level, edit the configuration file and restart the application.

#### Frontend Logging (JavaScript)

The application also includes structured logging in the browser console for all JavaScript operations:

- **AI Chat Interface**: Request/response tracking, error handling, UI state changes
- **Movie Board**: Form interactions, CRUD operations, modal management
- **TV Shows Board**: Form interactions, CRUD operations, modal management

Configure frontend logging using URL parameters or localStorage:

```javascript
// Set via URL parameter
http://localhost:8080/ai?log=DEBUG
http://localhost:8080/movie-board?log=WARN
http://localhost:8080/tv-shows-board?log=ERROR

// Or set programmatically in browser console
localStorage.setItem('ai_log_level', 'DEBUG');
localStorage.setItem('movie_log_level', 'INFO');
localStorage.setItem('tvshow_log_level', 'WARN');
```

**Available Log Levels:**
- `DEBUG`: Detailed debug information (form interactions, DOM updates)
- `INFO`: General information about user actions and operations
- `WARN`: Warning messages for potentially problematic situations
- `ERROR`: Error messages for failed operations

**Sample Frontend Log Output:**
```
[2024-07-25T17:11:17.123Z] INFO: AI chat interface initializing...
[2024-07-25T17:11:17.125Z] INFO: AI chat interface initialized successfully
[2024-07-25T17:11:20.456Z] INFO: Form submitted with message: What is Go?
[2024-07-25T17:11:20.458Z] INFO: Starting AI request for message: What is Go?
[2024-07-25T17:11:20.460Z] INFO: Sending request to /ai/query with message: What is Go?
[2024-07-25T17:11:22.789Z] DEBUG: Response received, status: 200
[2024-07-25T17:11:22.790Z] DEBUG: Raw response received, length: 1250
[2024-07-25T17:11:22.791Z] INFO: AI response displayed successfully
```

## Customization

### Adding New Routes

Add new handlers in `internal/handlers/handlers.go`:

```go
func NewRouteHandler(w http.ResponseWriter, r *http.Request) {
    // Your handler logic here
}
```

Then register the route in `internal/server/server.go`:

```go
http.HandleFunc("/new-route", handlers.NewRouteHandler)
```

### Styling

Modify the HTML template in `web/templates/index.html` and use Tailwind CSS classes. The application uses a fixed dark theme.

### HTMX Interactions

Add HTMX attributes to HTML elements:

```html
<button 
    hx-post="/api/endpoint" 
    hx-target="#target-element"
    hx-swap="innerHTML">
    Click me
</button>
```

To add new HTMX endpoints, create handlers in `internal/handlers/handlers.go` and register them in `internal/server/server.go`.

### Movie Board

The application includes a comprehensive movie management system:

- **Add Movies**: Form with title, year, genre, and notes
- **Movie List**: Display all added movies with delete functionality
- **HTMX Integration**: Real-time updates without page reloads
- **Responsive Design**: Works on all device sizes
- **In-Memory Storage**: Simple storage for demo purposes

To customize the movie board, modify the handlers in `internal/handlers/handlers.go`.

## Technologies Used

- **Go**: Backend web server and templating
- **HTMX**: Dynamic interactions and AJAX requests
- **Tailwind CSS**: Utility-first CSS framework
- **HTML5**: Semantic markup and modern web standards

## Browser Support

- Chrome/Chromium (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## License

This project is open source and available under the [MIT License](LICENSE).

## Contributing

Feel free to submit issues and enhancement requests!
