# URL Shortener

A full-stack URL shortener: a RESTful API built with Go, Gin, PostgreSQL, and Redis, with a React + TypeScript dashboard for managing URLs and viewing statistics. Fully containerized with Docker Compose, with JWT-based authentication for write operations.

## Features
- User registration and login with JWT authentication
- Create short URLs from long URLs (auth required)
- List all short URLs (auth required)
- Redirect to original URL using short code, with Redis caching for fast lookups
- Retrieve URL details
- Update existing short URLs (auth required)
- Delete short URLs (auth required)
- Track access statistics (hit count)
- Health check endpoint
- React + TypeScript dashboard: create, view, edit, delete URLs, view stats
- Fully containerized (API, Postgres, Redis, frontend) via Docker Compose
- Secrets/config managed via `.env` (not committed)
- CI via GitHub Actions (build + test on push/PR)

## Project Structure
```
shortURL/
├── api/
│   ├── Dockerfile
│   ├── main.go          # Application entry point, route setup
│   ├── handlers.go      # HTTP handlers (API logic)
│   ├── database.go      # DB + Redis initialization and queries
│   ├── models.go        # Data models
│   ├── auth.go          # JWT + password hashing logic
│   ├── middleware.go    # Auth middleware
│   └── utils.go         # Utility functions (short code generation)
├── frontend/
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── src/
│   │   ├── pages/       # Login, Register, Dashboard, URLDetail
│   │   ├── components/  # URLTable, CreateURLForm, ProtectedRoute
│   │   ├── api/         # Axios client + typed API calls
│   │   ├── context/     # AuthContext (JWT state)
│   │   └── types/       # TypeScript interfaces matching Go models
├── .github/workflows/   # CI pipeline
├── docker-compose.yml
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites
- Go (>= 1.26)
- Node.js (>= 24) — for local frontend dev
- Docker & Docker Compose

## Setup

### 1. Configure environment variables
```bash
cp .env.example .env
```
Fill in `.env` with your own values (DB credentials, JWT secret, etc.) — never commit this file.

### 2. Run with Docker Compose (Recommended)
```bash
git clone <repo-url>
cd shortURL
docker-compose up --build
```

- Frontend: `http://localhost`
- Backend API: `http://localhost:8080`

**Common Docker Compose commands:**
```bash
docker-compose up          # Start (reuse existing images)
docker-compose up --build  # Rebuild images and start
docker-compose stop        # Stop containers
docker-compose start       # Start stopped containers
docker-compose down        # Stop and remove containers
docker-compose down -v     # Stop, remove containers, and delete DB data
```

### 3. Run Locally (Without Docker)

**Backend:**
```bash
go mod tidy
$env:DATABASE_URL="host=localhost port=5432 user=postgres password=postgres dbname=urlshortener sslmode=disable"
$env:JWT_SECRET="your-secret-key"
$env:REDIS_ADDR="localhost:6379"
go run ./api
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```
Runs on `http://localhost:5173`

## API Endpoints

### Health Check
```
GET /health
```

### Register
```
POST /register
Body: { "username": "ashish", "password": "password123" }
Response: { "token": "..." }
```

### Login
```
POST /login
Body: { "username": "ashish", "password": "password123" }
Response: { "token": "..." }
```

### Create Short URL 🔒
```
POST /shorten
Headers: Authorization: Bearer <token>
Body: { "url": "https://example.com" }
```

### List All Short URLs 🔒
```
GET /shorten
Headers: Authorization: Bearer <token>
```

### Get Short URL Details
```
GET /shorten/:code
```

### Update Short URL 🔒
```
PUT /shorten/:code
Headers: Authorization: Bearer <token>
Body: { "url": "https://new-url.com" }
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

## Database & Caching

**PostgreSQL** with three tables:
- `users` — registered users and hashed passwords
- `urls` — short code to URL mappings
- `stats` — access history per short code

**Redis** caches short-code → URL lookups for fast redirects (cache-aside pattern), invalidated on update/delete.

Postgres data persists in a Docker volume (`postgres_data`) across container restarts.

## CI/CD

GitHub Actions workflow (`.github/workflows/ci.yml`) runs on every push/PR to `main`:
- Builds and tests the Go backend
- Builds the React frontend
- Verifies Docker images build successfully