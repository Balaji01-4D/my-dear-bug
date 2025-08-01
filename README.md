# 🐛 My Dear Bug

> *A modern, lightweight bug tracking and community-driven debugging platform built with Go and Gin*

[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Gin-Web%20Framework-green.svg)](https://gin-gonic.com)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue.svg)](https://postgresql.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Overview

**My Dear Bug** is a community-driven platform where developers can share, discover, and collaboratively solve coding bugs. Whether you're stuck on a tricky algorithm, encountering unexpected behavior, or want to help others debug their code, this platform provides a clean, efficient way to connect with the developer community.

### ✨ Key Features

- **Bug Discovery**: Browse and search through a curated collection of real-world coding bugs
- **Bug Submission**: Share your bugs with detailed descriptions and code snippets
- **Community Voting**: Upvote the most interesting or critical bugs to prioritize them
- **Language Filtering**: Filter bugs by programming language (Go, Python, JavaScript, etc.)
- **Trending Bugs**: Discover the most popular and trending bugs in the community
- **Admin Controls**: Secure administrative features for content moderation
- **High Performance**: Built with Go and Gin for optimal speed and efficiency
- **Robust Storage**: PostgreSQL database with GORM for reliable data persistence

## Architecture

### Clean Architecture Design
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controllers   │───▶│    Services     │───▶│  Repositories   │
│  (HTTP Layer)   │    │ (Business Logic)│    │  (Data Layer)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│      DTOs       │    │     Models      │    │   PostgreSQL    │
│ (Data Transfer) │    │  (Domain)       │    │   Database      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Project Structure
```
my-dear-bug/
├── config/                  # Configuration management
│   └── config.go           # Database & environment setup
├── internals/              # Core business logic
│   ├── bug/                # Bug management domain
│   │   ├── controller.go   # HTTP request handlers
│   │   ├── service.go      # Business logic
│   │   ├── repository.go   # Data access layer
│   │   ├── model.go        # Bug data model
│   │   └── dto.go          # Data transfer objects
│   ├── upvote/             # Voting system domain
│   │   ├── controller.go   # Upvote handlers
│   │   ├── service.go      # Voting logic
│   │   ├── repository.go   # Vote data access
│   │   └── model.go        # Upvote data model
│   └── middleware/         # HTTP middleware
│       └── adminAuth.go    # Admin authentication
├── migrate/                # Database migrations
│   ├── migrate.go          # Migration runner
│   └── init.sql            # Initial schema
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
└── Makefile                # Build automation
```

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Backend** | Go 1.24.2 | High-performance server runtime |
| **Web Framework** | Gin | Fast HTTP router and middleware |
| **Database** | PostgreSQL | Reliable relational database |
| **ORM** | GORM | Object-relational mapping |
| **Configuration** | godotenv | Environment variable management |
| **Build Tool** | Make | Build automation and task running |

## API Endpoints

### Bug Management
| Method | Endpoint | Description | Authentication |
|--------|----------|-------------|----------------|
| `GET` | `/bugs` | List all bugs with pagination | None |
| `GET` | `/bugs/{id}` | Get specific bug details | None |
| `POST` | `/bugs` | Submit a new bug report | None |
| `DELETE` | `/bugs/{id}` | Delete a bug | Admin only |

### Filtering & Discovery
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/bugs/language/{language}` | Filter bugs by programming language |
| `GET` | `/bugs/top` | Get highest-rated bugs |

### Community Voting
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/bugs/{id}/upvote` | Upvote a bug (IP-based deduplication) |

### Query Parameters
- `offset` - Skip number of records for pagination (default: 0)
- `limit` - Maximum records to return (default: 10, max: 100)

## 🚀 Quick Start

### Prerequisites
- Go 1.24.2 or later
- PostgreSQL database
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Balaji01-4D/my-dear-bug.git
   cd my-dear-bug
   ```

2. **Install dependencies**
   ```bash
   make install-deps
   ```

3. **Configure environment**
   
   Create a `.env` file in the project root:
   ```env
    DATABASE_URL=postgres://username:password@localhost:5432/mydearbug?sslmode=disable
    SERVER_ADDR=:8080

    # Admin Authentication
    ADMIN_USERNAME=admin
    ADMIN_PASSWORD=supersecret

   ```

4. **Set up database**
   ```bash
   # Run migrations
   go run migrate/migrate.go
   ```

5. **Start the application**
   ```bash
   make run
   ```

   🎉 **Your API is now running at:** `http://localhost:8080`

## Usage Examples

### Submit a Bug
```bash
curl -X POST http://localhost:8080/bugs \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Memory leak in Python list comprehension",
    "description": "Nested list comprehensions cause excessive memory usage",
    "snippet": "data = [[x*y for x in range(1000)] for y in range(1000)]",
    "language": "python"
  }'
```

### Browse Bugs
```bash
# Get all bugs with pagination
curl "http://localhost:8080/bugs?limit=10&offset=0"

# Filter by language
curl "http://localhost:8080/bugs/language/go?limit=5"

# Get trending bugs
curl "http://localhost:8080/bugs/top?limit=10"
```

### Vote on Bugs
```bash
# Upvote a bug
curl -X POST http://localhost:8080/bugs/1/upvote
```

### Sample Response
```json
{
  "id": 1,
  "title": "Deadlock Detected",
  "description": "Deadlock occurs when two goroutines wait on each other.",
  "language": "go",
  "snippet": "mutex1.Lock(); mutex2.Lock();",
  "createdAt": "2025-08-01T15:05:58.156094+05:30",
  "upvotes": 15
}
```

## Data Models

### Bug Entity
```go
type Bug struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Language    string    `json:"language"`
    Snippet     string    `json:"snippet"`
    CreatedAt   time.Time `json:"createdAt"`
    Upvotes     int       `json:"upvotes"`
}
```

### Upvote Entity
```go
type Upvote struct {
    ID        uint      `json:"id"`
    BugID     uint      `json:"bug_id"`
    IPHash    string    `json:"-"`        // Privacy-protected
    CreatedAt time.Time `json:"createdAt"`
}
```

## Security Features

- **IP-based Vote Deduplication**: Prevents spam voting using SHA-256 hashed IP addresses
- **Admin Authentication**: Secure middleware for administrative operations
- **Input Validation**: Comprehensive request validation using Gin's binding
- **SQL Injection Protection**: GORM ORM provides built-in protection

## Development

### Available Commands
```bash
make help          # Show all available commands
make run           # Start the development server
make build         # Build the application binary
make test          # Run the test suite
make clean         # Clean build artifacts
make fmt           # Format Go code
make lint          # Run code linting
```

### Development Workflow
1. Make your changes
2. Run `make fmt` to format code
3. Run `make test` to ensure tests pass
4. Run `make run` to test locally
5. Submit your pull request

## Contributing

We welcome contributions! Here's how you can help:

1. **Report Bugs**: Use the GitHub Issues tab
2. **Suggest Features**: Open a feature request
3. **Submit PRs**: Fork, branch, code, test, submit
4. **Improve Docs**: Help us make documentation better

### Contribution Guidelines
- Follow Go best practices and conventions
- Write tests for new features
- Update documentation for API changes
- Use meaningful commit messages

## Deployment

### Production Checklist
- [ ] Set `GIN_MODE=release` in production
- [ ] Use environment variables for sensitive data
- [ ] Enable database connection pooling
- [ ] Set up proper logging and monitoring
- [ ] Configure reverse proxy (Nginx/Apache)
- [ ] Enable HTTPS/TLS

### Environment Variables
```env
DATABASE_URL=          # PostgreSQL connection string
SERVER_ADDR=          # Server bind address (e.g., :8080)
GIN_MODE=             # Gin mode (debug/release)
```

## Performance

- **Fast Response Times**: Gin framework provides excellent performance
- **Efficient Queries**: GORM with optimized database queries
- **Pagination**: Built-in pagination for large datasets
- **Connection Pooling**: PostgreSQL connection pooling for scalability

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internals/bug/
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Email**: support@mydearbug.com
- **Issues**: [GitHub Issues](https://github.com/Balaji01-4D/my-dear-bug/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Balaji01-4D/my-dear-bug/discussions)

## Roadmap

- [ ] User authentication and profiles
- [ ] Email notifications for bug updates
- [ ] Advanced search with filters
- [ ] Mobile API endpoints
- [ ] Integration with GitHub/GitLab
- [ ] Analytics dashboard
- [ ] REST API documentation with Swagger
- [ ] Real-time updates with WebSockets

---

<div align="center">

**Built with ❤️ by the My Dear Bug community**

[⭐ Star this repo](https://github.com/Balaji01-4D/my-dear-bug) | [Report Bug](https://github.com/Balaji01-4D/my-dear-bug/issues) | [Request Feature](https://github.com/Balaji01-4D/my-dear-bug/issues)

</div>
