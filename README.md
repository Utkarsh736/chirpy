# Chirpy

Chirpy is a backend social media API built with Go. It serves as a Twitter-like platform where users can post short messages called "chirps" (limited to 140 characters), manage their accounts, and subscribe to premium membership features.

## Project Overview

This project was developed as part of the [Boot.dev](https://www.boot.dev) backend engineering curriculum. It demonstrates the implementation of a production-ready RESTful API with modern backend development practices including database management with PostgreSQL, secure authentication, and webhook integrations.

**Note**: This was a guided project completed through Boot.dev's HTTP Servers course. While following the curriculum, I implemented all core features and gained hands-on experience with professional backend development patterns.

## Features

### User Management
- **Account Creation**: Register new users with email and secure password hashing (Argon2id)
- **Authentication**: JWT-based access tokens (1-hour expiry) and refresh tokens (60-day expiry)
- **Profile Updates**: Change email and password for authenticated users
- **Token Management**: Refresh access tokens and revoke refresh tokens

### Chirps (Posts)
- **Create Chirps**: Post messages up to 140 characters with automatic profanity filtering
- **Retrieve Chirps**: Get all chirps or filter by author ID
- **Sorting**: Sort chirps by creation date (ascending or descending)
- **Delete Chirps**: Users can delete their own chirps with proper authorization checks
- **Profanity Filter**: Automatically replaces inappropriate words with `****`

### Premium Membership (Chirpy Red)
- **Webhook Integration**: Process payment provider (Polka) webhooks for membership upgrades
- **API Key Authentication**: Secure webhook endpoint with API key validation
- **Membership Status**: Track premium user status across all endpoints

### Security
- **Password Hashing**: Argon2id for secure password storage
- **JWT Authentication**: Stateless authentication with HS256 signing
- **API Key Protection**: Webhook endpoints secured with API keys
- **Authorization**: Resource ownership validation (users can only modify their own content)
- **HTTP Status Codes**: Proper 401 (Unauthorized) vs 403 (Forbidden) distinction

### Admin Features
- **Metrics Dashboard**: HTML-based admin page showing server statistics
- **Reset Endpoint**: Environment-gated endpoint to clear database (dev only)
- **Request Counter**: Middleware tracking fileserver hits

## Tech Stack

- **Language**: Go 1.22+ (using new routing enhancements)
- **Database**: PostgreSQL 15+
- **HTTP Router**: Standard library `net/http` with method-based routing
- **SQL Generation**: [SQLC](https://sqlc.dev) for type-safe SQL queries
- **Migrations**: [Goose](https://github.com/pressly/goose) for database schema management
- **Authentication**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt) for JWT handling
- **Password Hashing**: [argon2id](https://github.com/alexedwards/argon2id) library
- **Environment Config**: [godotenv](https://github.com/joho/godotenv) for local development

## API Endpoints

### Public Endpoints
- `GET /api/healthz` - Health check endpoint
- `POST /api/users` - Create new user account
- `POST /api/login` - Authenticate and receive tokens

### Authenticated Endpoints (Requires JWT)
- `PUT /api/users` - Update user email/password
- `POST /api/chirps` - Create a new chirp
- `DELETE /api/chirps/{chirpID}` - Delete own chirp
- `POST /api/refresh` - Get new access token using refresh token
- `POST /api/revoke` - Revoke a refresh token

### Read-Only Endpoints
- `GET /api/chirps` - Get all chirps (supports `?author_id=` and `?sort=asc|desc`)
- `GET /api/chirps/{chirpID}` - Get specific chirp by ID

### Webhook Endpoints
- `POST /api/polka/webhooks` - Handle payment provider webhooks (API key required)

### Admin Endpoints
- `GET /admin/metrics` - View server metrics (HTML dashboard)
- `POST /admin/reset` - Reset database (dev environment only)

### Static Assets
- `/app/*` - Fileserver for web interface

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 15 or higher
- `goose` CLI tool for migrations
- `sqlc` CLI tool for code generation

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/Utkarsh736/chirpy.git
   cd chirpy
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL**:
   ```bash
   # Start PostgreSQL service
   sudo service postgresql start
   
   # Create database
   sudo -u postgres psql
   CREATE DATABASE chirpy;
   \q
   ```

4. **Configure environment variables**:
   
   Create a `.env` file in the root directory:
   ```env
   DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
   PLATFORM=dev
   JWT_SECRET=<generate-with-openssl-rand-base64-64>
   POLKA_KEY=<insert_polka_key>
   ```

5. **Run database migrations**:
   ```bash
   cd sql/schema
   goose postgres "postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" up
   cd ../..
   ```

6. **Generate SQLC code** (if you modify queries):
   ```bash
   sqlc generate
   ```

7. **Build and run the server**:
   ```bash
   go build -o out && ./out
   ```

   The server will start on `http://localhost:8080`

### Testing the API

Example using `curl`:

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"securepass123"}'

# Login and get token
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"securepass123"}' \
  | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)

# Create a chirp
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"body":"Hello, Chirpy world!"}'

# Get all chirps (sorted newest first)
curl "http://localhost:8080/api/chirps?sort=desc"
```

## Project Structure

```
chirpy/
├── main.go                  # Main server application
├── .env                     # Environment variables (gitignored)
├── go.mod                   # Go module dependencies
├── sql/
│   ├── schema/              # Database migrations (Goose)
│   │   ├── 001_users.sql
│   │   ├── 002_chirps.sql
│   │   ├── 003_users_password.sql
│   │   ├── 004_refresh_tokens.sql
│   │   └── 005_users_chirpy_red.sql
│   └── queries/             # SQL queries (SQLC)
│       ├── users.sql
│       ├── chirps.sql
│       └── refresh_tokens.sql
├── internal/
│   ├── auth/                # Authentication helpers
│   │   ├── auth.go          # Password hashing, JWT, token extraction
│   │   └── auth_test.go     # Unit tests
│   └── database/            # Generated by SQLC
│       ├── db.go
│       ├── models.go
│       ├── users.sql.go
│       ├── chirps.sql.go
│       └── refresh_tokens.sql.go
├── assets/                  # Static assets
│   └── logo.png
└── index.html               # Homepage

```

## Key Learnings

Through building this project, I gained practical experience with:

- **RESTful API Design**: Implementing proper HTTP methods, status codes, and resource-based routing
- **Database Patterns**: Migrations, foreign keys, cascading deletes, and ACID transactions
- **Security Best Practices**: Password hashing, JWT authentication, API key validation, and authorization checks
- **Go Concurrency**: Using goroutines for handling multiple HTTP requests simultaneously
- **Type-Safe SQL**: Leveraging SQLC to generate Go code from SQL queries
- **Middleware Patterns**: Request counting and authentication middleware
- **Webhook Integration**: Processing external service callbacks with proper validation
- **Environment Configuration**: Separating dev/production configs with environment variables

## Future Enhancements

Potential features to add:
- Pagination for chirps endpoint
- User following/followers system
- Like/favorite functionality for chirps
- Rate limiting middleware
- Full-text search for chirps
- Image uploads
- WebSocket support for real-time updates

## Acknowledgments

This project was completed as part of Boot.dev's [Learn HTTP Servers](https://www.boot.dev/courses/learn-http-servers) course. Special thanks to Lane Wagner and the Boot.dev team for the excellent curriculum.

## License

This project is open source and available under the [MIT License](LICENSE).
