# Golang Boilerplate

Golang Boilerplate is your next step in web development.

## Setup

### Requirements

- [Go 1.20>](https://go.dev/dl/)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Postgres 15>](https://www.postgresql.org/download/)

### Running

Create a .env file inside configs folder.

> **NOTE:** You can copy and paste .EXAMPLE.env and rename it.

```sh
# .EXAMPLE.env
DATABASE_URL="mysql:mysql@tcp(127.0.0.1:3306)/database?parseTime=True"
PORT="4000"
ENABLE_STARTUP_MESSAGE=true
ENABLE_PRINT_ROUTES=false
ENABLE_STACK_TRACE=false
SKIP_MIGRATION=false
REDIS_URL="redis://:redis@127.0.0.1:6379"
EMAIL_SENDER="no-reply@golangboilerplate.com"
EMAIL_SENDER_NAME="Golang Boilerplate (No Reply)"
# ENVIRONMENT="DEV"
# SENDGRID_API_KEY="SG.my-send-grid-key"
```

Then you could start the dev environment dependencies, with docker compose command:

```sh
docker compose up -d
```

After that you could start the program:

```sh
go run main.go api
```

## Examples

### Login Request

```sh
# Request
curl -X POST \
  <BASE_URL>/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{
      "username": "admin",
      "password": "Temp#dev"
  }'

# Response
{
  "expiresAt": 1682430104,
  "refreshExpiresAt": "2023-05-01T10:41:44.278228-03:00",
  "refreshToken": "ad27ad8b-28b1-44ab-a793-af99c865aacc",
  "token": "<JwtToken>"
}
```

### Further

<!-- TODO: Create a docs subdomain like abp -->
See more at [docs.golangboilerplate.com](https://docs.golangboilerplate.com/)