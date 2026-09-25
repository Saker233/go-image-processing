# Go Image Processing Service

This project implements a backend service for uploading, listing, retrieving, deleting, and processing images. It follows the requirements from the Roadmap project: https://roadmap.sh/projects/image-processing-service

## Overview

The service provides:

- User registration and login
- JWT-based authentication for protected routes
- Image upload to S3-compatible storage
- Image metadata persistence in PostgreSQL
- Listing images for the authenticated user
- Fetching a single image by ID
- Deleting an image and its stored artifacts
- Image transformation endpoint placeholder for future processing logic

## Tech Stack

- Go
- Gin web framework
- PostgreSQL
- AWS S3 SDK
- JWT authentication
- SQLC generated database access

## Project Structure

- `main.go` - application entry point
- `internal/api/` - HTTP handlers and authentication middleware
- `internal/database/` - database connection, models, and generated SQL queries
- `internal/util/` - hashing, JWT helpers, and S3 utilities
- `database/migration/` - schema migration files
- `database/queries/` - SQL definitions used by sqlc

## Environment Configuration

Create or update `app.env` with values similar to:

```env
PORT=:8000
DB_DRIVER=postgres
DB_SOURCE=postgres://postgres:postgres@localhost:5433/image_processing?sslmode=disable
JWT_SECRET=your-secret-key
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-access-key
AWS_REGION=eu-central-1
```

## Run the Application

```bash
go mod download
make run
```

Or run directly:

```bash
go run .
```

## Database Setup

The project expects PostgreSQL to be running and the schema to be initialized with the migration in:

```bash
# example
psql -h localhost -p 5433 -U postgres -d image_processing -f database/migration/000001_init_schema.up.sql
```

## API Endpoints

### Authentication

#### Register

```http
POST /register
Content-Type: application/json
```

Request body:

```json
{
  "username": "user1",
  "password": "password123"
}
```

#### Login

```http
POST /login
Content-Type: application/json
```

Request body:

```json
{
  "username": "user1",
  "password": "password123"
}
```

### Image Routes

All image routes require a bearer token:

```http
Authorization: Bearer <token>
```

#### Upload Image

```http
POST /images
Content-Type: multipart/form-data
```

Form field:

- `file`: image file

#### List Images

```http
GET /images?page=1&limit=10
```

#### Get Image

```http
GET /images/:id
```

#### Delete Image

```http
DELETE /images/:id
```

#### Transform Image

```http
POST /images/:id/transform
Content-Type: application/json
```

Example body:

```json
{
  "transformations": {
    "resize": { "width": 800, "height": 600 },
    "format": "jpeg"
  }
}
```

## Image Processing Notes

This service currently covers the core backend flow and infrastructure for image storage and metadata management. The transformation route is present and structured for adding actual image processing logic such as:

- resize
- crop
- rotate
- flip and mirror
- grayscale/sepia filters
- format conversion

## Roadmap Reference

This project is based on the backend challenge from:

https://roadmap.sh/projects/image-processing-service

## License

This project is for learning and project development purposes.
