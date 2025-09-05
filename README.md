# 🐛 My Dear Bug

> A modern, lightweight community-driven debugging platform built with Go and Gin

[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Gin-Web%20Framework-green.svg)](https://gin-gonic.com)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue.svg)](https://postgresql.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Overview

My Dear Bug is a community platform where developers share real-world issues ("confessions"), discover tricky bugs, and help each other debug. Submit a confession with details and code, upvote interesting ones, and explore by language or tags.

## Architecture

Clean layered architecture with clear separation of concerns:

```
Controllers → Services → Repositories → Database (PostgreSQL)
            DTOs ↔ Models
```

### Project Structure
```
shit-happens/
├── config/                  # Configuration & DB setup
│   ├── config.go            # Loads env (DATABASE_URL, PORT) & InitDB
│   └── devConfig.go         # .env loader (godotenv) for local dev
├── internals/
│   ├── confession/          # Confession (bug report) domain
│   │   ├── controller.go    # HTTP handlers (/confessions)
│   │   ├── service.go       # Business logic
│   │   ├── repository.go    # GORM queries (pagination, filters)
│   │   ├── model.go         # Confession entity + relations
│   │   └── dto.go           # Request validation
│   ├── upvote/              # Upvote domain
│   │   ├── controller.go    # POST /confessions/:id/upvote
│   │   ├── service.go       # Upvote logic (+1 counter)
│   │   ├── repository.go    # HasUpvoted + persist
│   │   └── model.go         # Upvote entity
│   ├── tag/                 # Tagging & suggestions
│   │   ├── controller.go    # /tags endpoints
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── model.go         # Tag entity
│   └── middleware/
│       ├── adminAuth.go     # Basic auth for protected routes
│       └── rateLimit.go     # Per-IP POST rate limiting
├── migrate/
│   ├── migrate.go           # GORM AutoMigrate
│   └── init.sql             # (Legacy) SQL schema
├── main.go                  # App entrypoint & route registration
├── Makefile                 # build/run/test/fmt tasks
├── go.mod / go.sum
└── README.md
```

## Technology Stack

- Go + Gin web framework
- PostgreSQL + GORM ORM
- godotenv (local dev)
- Make for developer tasks

## API Endpoints

### Confession Management
- GET  `/confessions` — List with pagination (offset, limit)
- GET  `/confessions/:id` — Get details
- POST `/confessions` — Create a confession (rate-limited per IP)
- DELETE `/confessions/:id` — Delete (admin only)

### Filtering & Discovery
- GET `/confessions/language/:language` — Filter by language (case-insensitive)
- GET `/confessions/top` — Highest upvoted confessions
- GET `/confessions/trending/weekly` — Trending confessions for the last 7 days
- GET `/confessions/trending/monthly` — Trending confessions for the last 30 days
- GET `/confessions/hall-of-fame` — All-time notable (e.g. high-impact) confessions
- GET `/confessions/random` — Random selection (use for inspiration / shuffle)

### Community Voting
- POST `/confessions/:id/upvote` — Upvote (deduplicated by IP hash)

### Tags
- GET `/tags` — List tags
- POST `/tags` — Create tag
- GET `/tags/suggest?query=<prefix>` — Autocomplete suggestions
- DELETE `/tags/:id` — Delete tag (admin only)

### Query Parameters
- `offset` — Pagination offset (default: 0 if omitted)
- `limit` — Pagination limit (if omitted, no limit is applied by the DB layer)

## Usage Examples

### Submit a Confession
```bash
curl -X POST http://localhost:8080/confessions \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Deadlock Detected",
    "description": "Two goroutines wait on each other causing deadlock",
    "snippet": "mutex1.Lock(); mutex2.Lock();",
    "language": "go",
    "tags": ["concurrency", "mutex"]
  }'
```

### Browse Confessions
```bash
# Get all with pagination
curl "http://localhost:8080/confessions?limit=10&offset=0"

# Filter by language
curl "http://localhost:8080/confessions/language/go?limit=5"

# Highest upvoted
curl "http://localhost:8080/confessions/top?limit=10"

# Trending this week
curl "http://localhost:8080/confessions/trending/weekly?limit=10"

# Trending this month
curl "http://localhost:8080/confessions/trending/monthly?limit=10"

# Hall of fame
curl "http://localhost:8080/confessions/hall-of-fame?limit=10"

# Random selection
curl "http://localhost:8080/confessions/random"
```

### Vote on a Confession
```bash
curl -X POST http://localhost:8080/confessions/1/upvote
```

### Tag APIs
```bash
# List tags
curl http://localhost:8080/tags

# Create tag
curl -X POST http://localhost:8080/tags -H 'Content-Type: application/json' -d '{"name": "performance"}'

# Suggest tags (prefix search)
curl "http://localhost:8080/tags/suggest?query=per"
```

### Sample Confession Response
```json
{
  "id": 1,
  "title": "Deadlock Detected",
  "description": "Two goroutines wait on each other.",
  "language": "go",
  "snippet": "mutex1.Lock(); mutex2.Lock();",
  "tags": [{ "id": 3, "name": "concurrency" }],
  "sentiment": "happy",
  "isFlagged": false,
  "createdAt": "2025-08-01T15:05:58.156094+05:30",
  "upvotes": 15
}
```

## Data Models

### Confession
```go
type Confession struct {
    ID          uint       `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Language    string     `json:"language"`
    Snippet     string     `json:"snippet"`
    Tags        []tag.Tag  `json:"tags"`           // many2many: confession_tags
    Sentiment   string     `json:"sentiment"`
    IsFlagged   bool       `json:"isFlagged"`
    CreatedAt   time.Time  `json:"createdAt"`
    Upvotes     int        `json:"upvotes"`
}
```

### Upvote
```go
type Upvote struct {
    ID           uint      `json:"id"`
    ConfessionID uint      `json:"confessionId"`
    IPHash       string    `json:"-"`           // privacy protected
    CreatedAt    time.Time `json:"createdAt"`
}
```

### Tag
```go
type Tag struct {
    ID   uint   `json:"id"`
    Name string `json:"name"` // unique
}
```

## Security & Middleware

- IP-based upvote deduplication (SHA-256 hash of client IP)
- Admin authentication (Basic Auth) for DELETE endpoints
- Rate limiting: 10 POSTs/hour per IP (burst 3) on confession creation

## Development

### Prerequisites
- Go 1.24.2+
- PostgreSQL
- Git

### Install & Run
1) Install dependencies
```bash
make install-deps
```

2) Configure environment (.env)
```env
DATABASE_URL=postgres://username:password@localhost:5432/shithappens?sslmode=disable
PORT=:8080
ADMIN_USERNAME=admin
ADMIN_PASSWORD=supersecret
```

3) Run migrations
```bash
go run migrate/migrate.go
```

4) Start the API
```bash
make run
```
API: http://localhost:8080

### Available Commands
```bash
build                Build the application binary
fmt                  Format Go code
help                 Show all available commands
install-deps         Install all Go dependencies
run                  Start the development server
test                 Run the test suite
```

### Development Workflow
1. Make changes
2. make fmt
3. make test
4. make run
5. Open a PR

## Configuration Notes

- The server address is read from `PORT` (e.g., `:8080`)
- Gin is started in release mode by default in `main.go`
- Migrations use GORM AutoMigrate for `Confession` and `Upvote` (and will create tag relations)

## Testing

```bash
make test
# with coverage
go test -cover ./...
```

## Deployment

- Set `PORT` appropriately for your environment
- Use environment variables for secrets (no hardcoding)
- Configure reverse proxy and TLS as needed

## Roadmap

- User authentication & profiles
- Email notifications
- Advanced search & filters
- Mobile-friendly endpoints
- Swagger/OpenAPI docs
- Real-time updates (WebSockets)

---

<div align="center">
Built with ❤️ by the My Dear Bug community
</div>
