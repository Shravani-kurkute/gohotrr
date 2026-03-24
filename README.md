# GoLiveReload

A lightweight live reload server built in Go that automatically refreshes your browser when files change.

## Features

- **File Watching**: Monitors file changes in the current directory
- **WebSocket Communication**: Real-time notifications to connected clients
- **Static File Serving**: Serves static files from the current directory
- **Configurable Port**: Specify custom port via command-line argument
- **Simple & Fast**: Minimal dependencies, fast startup time

## Requirements

- Go 1.24.0 or higher

## Installation

1. Clone the repository or download the source code
2. Navigate to the project directory
3. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

### Running the Server

Start the server with the default port (8080):
```bash
go run main.go
```

Or specify a custom port:
```bash
go run main.go 3000
```

The server will start and output:
```
Server started at http://localhost:8080
```

### Building

To build an executable:
```bash
go build -o gohotrr .
```

Then run it:
```bash
./gohotrr
```

Or with a custom port:
```bash
./gohotrr 3000
```

## Dependencies

- **fsnotify**: For monitoring file system events
- **gorilla/websocket**: For WebSocket connections

## How It Works

1. The server watches the current directory for file changes using `fsnotify`
2. When a file is modified, the server notifies all connected WebSocket clients
3. Clients receive the notification and automatically reload the page
4. The server also injects reload middleware into served HTML files

## Architecture

- **File Watcher**: Watches for changes in the current directory recursively
- **WebSocket Handler**: Manages WebSocket connections for real-time updates
- **Reload Middleware**: Injects reload logic into HTTP responses
- **File Server**: Serves static files from the current directory

## License

MIT
