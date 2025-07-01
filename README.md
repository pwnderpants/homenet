# HTMX + Go + Tailwind CSS Boilerplate

A modern web application boilerplate built with HTMX for dynamic interactions, Go for the backend, and Tailwind CSS for styling with dark mode support.

## Features

- 🚀 **HTMX Integration**: Dynamic interactions without writing JavaScript
- ⚡ **Go Backend**: Fast, simple, and reliable web server
- 🎨 **Tailwind CSS**: Utility-first CSS framework
- 🌙 **Dark Theme**: Fixed dark theme design
- 📱 **Responsive Design**: Mobile-first responsive layout
- 🔄 **Interactive Counter**: Demo of HTMX functionality

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go      # Application entry point
├── internal/
│   ├── handlers/
│   │   └── handlers.go  # HTTP request handlers
│   ├── server/
│   │   └── server.go    # Server configuration
│   └── templates/
│       └── templates.go # Template management
├── web/
│   ├── templates/
│   │   └── index.html   # HTML templates
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
go build -o htmx-app cmd/server/main.go

# Run the binary
./htmx-app
```

You can also set a custom port using the `PORT` environment variable:

```bash
PORT=3000 go run cmd/server/main.go
```

## Features Explained

### HTMX Integration

The application demonstrates HTMX functionality with an interactive counter:

- Click the "Increment" button to see HTMX in action
- The counter updates without a full page reload
- Server-side rendering with Go templates

### Dark Theme

- **Fixed Dark Design**: Application uses a consistent dark theme
- **Optimized Colors**: Carefully selected dark colors for optimal readability
- **Custom Scrollbar**: Dark-themed scrollbar that matches the overall design

### Tailwind CSS

- **Utility Classes**: Rapid UI development with utility classes
- **Dark Mode Variants**: Automatic dark mode styling with `dark:` prefix
- **Responsive Design**: Mobile-first responsive breakpoints
- **Custom Styling**: Custom scrollbar styling for dark mode

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
