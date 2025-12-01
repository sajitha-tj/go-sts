# go-sts

`go-sts` is a Go-based implementation of a Security Token Service (STS) that supports OAuth2 flows. It is implemented as a PoC (Proof of Concept) and the service is built on top of the [Fosite](https://github.com/ory/fosite) library, which provides a robust framework for OAuth2 and OpenID Connect.

## Grant Types

The service supports the following OAuth2 grant types:
- **Authorization Code Grant**: Used for server-side applications where the client secret is kept confidential.
- **Client Credentials Grant**: Used for machine-to-machine communication where the client is also the resource owner.

## Project Structure

```
go-sts/
├── cmd/                             # Main entry point for the application
│   └── go-sts/
│       └── main.go                  # Application entry point
├── internal/
│   ├── app/
│   │   └── app.go                   # Application initialization and configuration
│   ├── configs/                     # Configuration management, including constants and loading environment variables
│   ├── lib/                         # Utility libraries
│   ├── middleware/                  # Middleware for HTTP requests
│   │   └── ctxMiddleware.go         # Context middleware
│   ├── repository/
│   │   ├── client_repository/       # Client repository logic
│   │   ├── issuer_repository/       # Issuer repository logic
│   │   ├── session_repository/      # Session repository logic
│   │   └── user_repository/         # User repository logic
│   ├── routes/
│   │   └── oauth_routes.go          # OAuth2 related routes
│   ├── service/
│   │   ├── authentication_service/  # User authentication related services
│   │   └── oauth_provider/          # Fosite OAuth2 provider setup
│   └── storage/                     # Fosite's storage interface implementation
│   └── templates/                   # HTML templates for rendering
├── resources/                       # Resource files
├── setup/
│   └── testDB.go                    # Test database setup
├── client/                          # OAuth2 client for testing
├── docker-compose.yaml              # Docker Compose configuration
├── go.mod                           # Go module file
├── .env                             # Environment variables file
└── README.md                        # Project documentation
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
   PORT=8080
   DB_USERNAME=oauth
   DB_PASSWORD_FILE=/home/sajithaj/my-sts-project/go-sts/resources/db_password
   DB_NAME=oauthdb
   FOSITE_SECRET_FILE=/home/sajithaj/my-sts-project/go-sts/resources/sign_secret
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

### Auth Code Flow

1. Open your browser and navigate to the [`http://localhost:3846`](http://localhost:3846) URL. You can start the Authorization Code flow by clicking the "Authorize" button.
2. You will be redirected to the authorization server, where you can log in and authorize the client application. Use the following credentials:
   - **Username**: `peter`
   - **Password**: `secret`
3. After successful authentication, you will be redirected back to the client server with an authorization code.
4. Click on "Get Token" to exchange the authorization code for an access token.

### Client Credentials Flow

Use the following curl command to test the Client Credentials flow:

```bash
curl --location '123e4567-e89b-12d3-a456-426614174000.localhost:8080/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=client_credentials' \
--data-urlencode 'client_id=my-client' \
--data-urlencode 'client_secret=foobar'
```
