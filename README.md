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
├── go.mod                  # Go module file
├── go.sum                  # Go module dependencies
└── README.md               # Project documentation
```

## Prerequisites

To run the project, you need to have the following installed:

- [Go 1.18 or later](https://golang.org/dl/)
- [PostgreSQL](https://www.postgresql.org/download/) (for local development)
- [Docker](https://www.docker.com/get-started) (for running PostgreSQL in a container)

## Getting Started

1. Clone the repository:

   ```bash
   git clone github.com/sajitha-tj/go-sts.git
   cd go-sts
   ```

2. Set up the PostgreSQL database:
   - You can either set up PostgreSQL locally or use Docker to run a PostgreSQL container.

   ```bash
   docker compose up -d
   ```
    - This will start a PostgreSQL container with the database `go_sts` and user `postgres` with password `password`.

3. Set up the environment variables:
   - Create a `.env` file in the root directory and add the following variables:

   ```plaintext
   DB_USER=oauth
   DB_PASSWORD=secret
   DB_NAME=oauthdb

   PORT=8080
   SIGNING_SECRET_FILE_PATH=/home/sajithaj/my-sts-project/go-sts/sign_secret.txt
   ```

4. Install the dependencies:

   ```bash
   go mod download
   ```

5. Run the client server to handle redirects:

    ```bash
    go run ./setup/clientServer/main.go
    ```
    This will start a simple HTTP server on port `3846` to handle the redirect from the authorization server.

6. Run the application:
    
    ```bash
    go run ./cmd/go-sts/main.go
    ```

## Testing

### Authorize request

Following request will return the auth code to the redirect_uri.

```bash
curl --location 'localhost:8080/authorize' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'response_type=code' \
--data-urlencode 'client_id=my-client' \
--data-urlencode 'redirect_uri=http://localhost:3846/callback' \
--data-urlencode 'scope=fosite openid photos offline' \
--data-urlencode 'state=random-state-value' \
--data-urlencode 'nonce=random-nonce-value' \
--data-urlencode 'code_challenge=example-code-challenge' \
--data-urlencode 'code_challenge_method=S256' \
--data-urlencode 'username=peter' \
--data-urlencode 'password=secret'
```

### Token request
