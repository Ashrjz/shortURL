# ShortURL

This is a solution to the URL shortener project on roadmap.sh : https://roadmap.sh/projects/url-shortening-service

**Built as a learning project.**

## Features
- Create short URLs from long URLs
- Redirect to original URL using short code
- Retrieve URL details
- Update existing short URLs
- Delete short URLs
- Track access statistics (hit count)
- Health check endpoint

## Project Structure
**api**<br>
&emsp;├── main.go # Application entry point, route setup <br>
&emsp;├── handlers.go # HTTP handlers (API logic) <br>
&emsp;├── database.go # DB initialization and queries <br>
&emsp;├── models.go # Data models <br>
&emsp;├── utils.go # Utility functions (short code generation) <br>
urls.db # SQLite database (auto-created) <br>


## Prerequisites
- Go (>= 1.18)
- SQLite

## Setup
```bash
# Clone repository
git clone <repo-url>
cd <project-folder>

# Install dependencies
go mod tidy

# Run application
go run main.go

Server runs on:
http://localhost:8080
```

## API Endpoints
Health Check
```
GET /health
```

Create Short URL
```
POST /shorten
Body:
{
  "url": "https://example.com"
}
```

Get Short URL Details
```
GET /shorten/:code
```

Update Short URL
```
PUT /shorten/:code
Body:
{
  "url": "https://new-url.com"
}
```

Delete Short URL
```
DELETE /shorten/:code
```

Get URL Statistics
```
GET /shorten/:code/stats
```

Redirect to Original URL
```
GET /:code
```