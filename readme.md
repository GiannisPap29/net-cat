# Net-Cat: TCP Chat Server

A Go implementation of a TCP-based chat server inspired by the NetCat utility, featuring a multi-client architecture with real-time messaging capabilities.

## Features

- TCP connection supporting multiple simultaneous clients (1-to-many relationship)
- User authentication with unique usernames
- Maximum connection limit (10 clients)
- Real-time message broadcasting
- Timestamp-based message formatting
- Chat history for new clients
- Join/Leave notifications
- Empty message filtering
- Configurable port settings (default: 8989)

## Requirements

- Go 1.23.4 or higher

## Installation

1. Clone the repository:
```bash
git clone [your-repository-url]
cd net-cat
```

2. Build the project:
```bash
go build
```

## Usage

### Starting the Server

- With default port (8989):
```bash
go run .
```

- With custom port:
```bash
go run . 2525
```

### Connecting to the Server

You can connect to the server using the `nc` (netcat) command:
```bash
nc localhost 8989
```

Upon connection, you will see:
1. Welcome banner with ASCII art
2. Name prompt
3. Chat history (if any)
4. Real-time messages from other users

### Message Format

Messages in the chat are formatted as follows:
```
[YYYY-MM-DD HH:MM:SS][username]: message
```

## Project Structure

```
.
├── internal/
│   ├── config/
│   │   └── config.go      # Configuration constants
│   ├── server/
│   │   ├── client.go      # Client handling
│   │   └── server.go      # Server implementation
│   └── utils/
│       └── utils.go       # Utility functions
├── main.go                # Entry point
└── go.mod                 # Go module file
```

## Implementation Details

- Uses Go routines for concurrent client handling
- Implements mutex locks for thread-safe operations
- Channel-based communication between clients and server
- Clean disconnection handling
- Robust error management

## Error Handling

The server handles various error scenarios:
- Invalid port numbers
- Connection failures
- Duplicate usernames
- Server capacity limits
- Network disconnections

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the [Your License] - see the LICENSE file for details