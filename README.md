# Music API App

## Overview

This application provides a REST API interface for accessing music artist information using the Gin framework in Go. The app serves as a middleware to query music data, including artist searches, artist details, and top tracks.

## Features

- Artist search functionality
- Retrieval of artist details
- Access to artist's top tracks by region

## Prerequisites

- Go 1.17 or later
- Gin framework
- Access to music API credentials

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/music-api-app.git
   cd music-api-app
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up environment variables:
   ```
   export API_CLIENT_ID=your_client_id
   export API_CLIENT_SECRET=your_client_secret
   ```

## Configuration

Update the `config.go` file with your API credentials and settings:

```go
// Example configuration
type Config struct {
    ClientID     string
    ClientSecret string
    RedirectURI  string
    APIBaseURL   string
}
```

For security reasons, update the trusted proxies configuration:

```go
router.SetTrustedProxies([]string{"127.0.0.1"})
```

## Usage

1. Start the server:
   ```
   go run main.go
   ```

2. The API will be available at `http://localhost:8080`

3. API Endpoints:
   - Search for an artist: `GET /api/search?q=artist_name&type=artist`
   - Get artist details: `GET /api/artist/{artist_id}`
   - Get artist's top tracks: `GET /api/artist/{artist_id}/top-tracks?market={country_code}`

## Common Issues

- `400 Bad Request` when accessing artist details: Ensure you're using a valid artist ID obtained from search results
- `502 Bad Gateway` during search: Check your API credentials and network connection
- Proxy warning: Configure trusted proxies as recommended in the Gin documentation

## License

[Your chosen license]

## Support

For questions or issues, please open an issue in the GitHub repository or contact [your contact information].