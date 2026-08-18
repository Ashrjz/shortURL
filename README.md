# URL Shortener API
 
A simple RESTful API for shortening URLs, built with Go, Gin, and PostgreSQL. Containerized with Docker and Docker Compose, with JWT-based authentication for write operations.

## Features
- User registration and login with JWT authentication
- Create short URLs from long URLs (auth required)
- Redirect to original URL using short code
- Retrieve URL details
- Update existing short URLs (auth required)
- Delete short URLs (auth required)
- Track access statistics (hit count)
- Health check endpoint
- Dockerized app + PostgreSQL via Docker Compose
## Project Structure
```
shortURL/
├── api/
│   ├── main.go         # Application entry point, route setup
│   ├── handlers.go     # HTTP handlers (API logic)
│   ├── database.go     # DB initialization and queries
│   ├── models.go       # Data models
│   ├── auth.go         # JWT + password hashing logic
│   ├── middleware.go   # Auth middleware
│   └── utils.go        # Utility functions (short code generation)
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```
 
## Prerequisites
- Go (>= 1.21)
- Docker & Docker Compose
## Setup
 
### Run with Docker Compose (Recommended)
```bash
# Clone repository
git clone <repo-url>
cd url-shortener
 
# Build and start app + PostgreSQL
docker-compose up --build
```
 
Server runs on: `http://localhost:8080`
 
**Common Docker Compose commands:**
```bash
docker-compose up          # Start (reuse existing image)
docker-compose up --build  # Rebuild image and start
docker-compose stop        # Stop containers
docker-compose start       # Start stopped containers
docker-compose down        # Stop and remove containers
docker-compose down -v     # Stop, remove containers, and delete DB data
```
 
### Run Locally (Without Docker)
```bash
# Install dependencies
go mod tidy
 
# Set environment variables
$env:DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=urlshortener sslmode=disable"
$env:JWT_SECRET="your-secret-key"
 
# Run application
go run ./api
```
 
## API Endpoints
 
### Health Check
```
GET /health
```
 
### Register
```
POST /register
Body:
{
  "username": "ashish",
  "password": "password123"
}
Response: { "token": "..." }
```
 
### Login
```
POST /login
Body:
{
  "username": "ashish",
  "password": "password123"
}
Response: { "token": "..." }
```
 
### Create Short URL 🔒
```
POST /shorten
Headers: Authorization: Bearer <token>
Body:
{
  "url": "https://example.com"
}
```
 
### Get Short URL Details
```
GET /shorten/:code
```
 
### Update Short URL 🔒
```
PUT /shorten/:code
Headers: Authorization: Bearer <token>
Body:
{
  "url": "https://new-url.com"
}
```
 
### Delete Short URL 🔒
```
DELETE /shorten/:code
Headers: Authorization: Bearer <token>
```
 
### Get URL Statistics
```
GET /shorten/:code/stats
```
 
### Redirect to Original URL
```
GET /:code
```
 
🔒 = Requires `Authorization: Bearer <token>` header
 
## Database
 
PostgreSQL with three tables:
- `users` — registered users and hashed passwords
- `urls` — short code to URL mappings
- `stats` — access history per short code
Data persists in a Docker volume (`postgres_data`) across container restarts.