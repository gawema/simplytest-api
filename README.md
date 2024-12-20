# Medications API

A RESTful API built with Go for managing medications. This API provides CRUD operations for medication records and uses MongoDB for data storage.

## Features

- RESTful endpoints for medications
- MongoDB integration
- Graceful shutdown
- Environment-based configuration
- Structured project layout

## Project Structure

The project follows a clean architecture pattern:

- `api/`: Contains all HTTP-related code
  - `handlers/`: HTTP request handlers
  - `routes.go`: Route definitions
- `storage/`: Contains all data-related code
  - `mongodb/`: MongoDB connection and operations
  - `models/`: Data structures
- `main.go`: Application entry point
- `.env`: Environment configuration

## Prerequisites

- Go 1.16 or higher
- MongoDB
- gorilla/mux package
- mongo-driver package

## Environment Setup

Required environment variables in `.env`:

- MONGODB_URI: MongoDB connection string
- MONGODB_DB_NAME: Database name
- MONGODB_COLLECTION: Collection name

## Installation

1. Clone the repository
2. Install dependencies with `go mod download`
3. Create `.env` file with required variables
4. Run with `go run .`

The server will start on `localhost:8080`

## API Documentation

### Medication Model

A medication object has the following structure:
- id: string (MongoDB ObjectID)
- name: string
- description: string

### Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /medications | List all medications |
| GET | /medications/{id} | Get single medication |
| POST | /medications | Create medication |
| PUT | /medications/{id} | Update medication |
| DELETE | /medications/{id} | Delete medication |

### Example Usage

GET all medications:
`curl http://localhost:8080/medications`

GET single medication:
`curl http://localhost:8080/medications/{id}`

CREATE medication:
`curl -X POST http://localhost:8080/medications -H "Content-Type: application/json" -d '{"name":"Ibuprofen","description":"Pain reliever"}'`

UPDATE medication:
`curl -X PUT http://localhost:8080/medications/{id} -H "Content-Type: application/json" -d '{"name":"Updated Name","description":"Updated description"}'`

DELETE medication:
`curl -X DELETE http://localhost:8080/medications/{id}`

## Live Demo

The API is deployed at: https://simplytest-api.onrender.com

**Note:** Since I am running on a free tier, the instance will spin down with inactivity, which can delay requests by 50 seconds or more.

Example:
```bash
curl https://simplytest-api.onrender.com/medications
```

