# URL Shortener Service

A simple URL shortener REST API built with Go.

## Features

- Shorten URLs with a REST API
- Redirect from short code to original URL
- Return same short code for duplicate URLs
- View top 3 most-shortened domains
- In-memory storage
- Thread-safe operations
- Comprehensive unit tests

## Project Structure

```
urlshortener/
├── main.go
├── go.mod
├── internal/
│   ├── handler/
│   │   └── handler.go
│   ├── service/
│   │   ├── url_service.go
│   │   └── url_service_test.go
│   └── storage/
│       ├── store.go
│       └── store_test.go
├── Dockerfile
└── README.md
```

## Building and Running

### Local

```bash
go build -o urlshortener .
./urlshortener
```

### Docker

```bash
docker build -t urlshortener:latest .
docker run -p 8080:8080 urlshortener:latest
```

## API Endpoints

### 1. Shorten URL

**Request:**
```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.example.com/very/long/path"}'
```

**Response:**
```json
{
  "short_code": "abc1234",
  "url": "https://www.example.com/very/long/path"
}
```

### 2. Redirect to Original URL

**Request:**
```bash
curl -L http://localhost:8080/redirect/abc1234
```

This will redirect to the original URL.

### 3. Get Top Domains Metrics

**Request:**
```bash
curl http://localhost:8080/metrics/top-domains
```

**Response:**
```json
{
  "top_domains": [
    {
      "domain": "udemy.com",
      "count": 6
    },
    {
      "domain": "youtube.com",
      "count": 4
    },
    {
      "domain": "wikipedia.org",
      "count": 2
    }
  ]
}
```

## Running Tests

```bash
go test ./...
```

For verbose output:
```bash
go test -v ./...
```

## Architecture

The application follows a clean three-layer architecture:

1. **Handler Layer**: Handles HTTP requests/responses
2. **Service Layer**: Contains business logic for URL shortening and metrics
3. **Storage Layer**: Manages data persistence (in-memory)

This separati