# go-sts

`go-sts` is a Go-based implementation of a Security Token Service (STS) that supports OAuth2 flows. It is implemented as a PoC (Proof of Concept) and the service is built on top of the [Fosite](github.com/ory/fosite) library, which provides a robust framework for OAuth2 and OpenID Connect.

## Grant Types

The service supports the following OAuth2 grant types:
- **Authorization Code Grant**: Used for server-side applications where the client secret is kept confidential.

## Project Structure

```plaintext
go-sts/
├── cmd/                     # Main entry point for the application
│   └── go-sts/
│       └── main.go          # Application startup logic
├── internal/
│   ├── controller/          # HTTP controllers for handling API endpoints
│   ├── lib/                 # Utility libraries (e.g., JSON handling)
│   ├── repository/          # Database models and data access logic
│   ├── service/             # Business logic and OAuth2 service implementation
│   ├── storage/             # Database connection and storage management
├── setup/                   # Setup scripts for testing and database initialization
│   ├──clientServer/         # Sample client to help handle redirects
│   ├── testDB.go            # Database setup for testing
├──                          # Go module dependencies
├── docker-compose.yaml      # Docker configuration for PostgreSQL
