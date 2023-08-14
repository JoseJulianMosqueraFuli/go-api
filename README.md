# Simple API using Golang and Gin Framework

This is a simple API built using the Golang programming language and the Gin web framework. The API provides an endpoint for creating and managing deliveries.

## Requirements

- Docker Desktop (activated)
- Golang

## Getting Started

1. Enter to the project directory:

```bash
git clone git@github.com:JoseJulianMosqueraFuli/go-api.git
```

2. Enter to the project directory:

```bash
cd go-api
```

3. Initialize the Go module:

```bash
go mod init myapi
```

4. Install the required dependencies:

```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/jinzhu/gorm
go get -u github.com/jinzhu/gorm/dialects/postgres
go get -u github.com/google/uuid
go get -u github.com/dgrijalva/jwt-go
```

5. Build the project:

```bash
docker build
```

6. Run the container:

```bash
docker up
```

The API will be accessible at http://localhost:8080/.

## API Endpoints: Available endpoints

### POST /register: Create a new registration user:

```bash
{
  'username' : '<username>',
  'password' : '<password>'
}
```

### POST /login: Access user token:

```bash
{
  'username' : '<username>',
  'password' : '<password>'

}

```

### POST /deliveries: Create a new delivery.

- Example JSON request body:

```json
{
  "id": "1234567890",
  "state": "pending",
  "pickup": {
    "pickup_lat": 37.7749,
    "pickup_lon": -122.4194
  },
  "dropoff": {
    "dropoff_lat": 34.0522,
    "dropoff_lon": -118.2437
  },
  "zone_id": "zone123"
}
```

### GET /deliveries/:id :

```bash

```

### GET /deliveries/creator/:creatorID :

```bash

```

### GET /deliveries:

```bash

```

### POST /deliveries/by-date: :

```bash

```

### POST /bots:

```bash

```

## Author

Build it by [Jose Julian Mosquera Fuli](https://github.com/JoseJulianMosqueraFuli).
