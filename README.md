# ProjectHub
ProjectHub is a backend-driven web application that allows authenticated users to submit projects, optionally upload images, and have them reviewed and approved by administrators before being published publicly.

It is designed as a production-oriented Go backend with a clean architecture, secure file handling, and a moderation workflow.

**Observability Focus:** The application implements structured logging from the start. See more about the approach in the "Structured Logging & Observability" section. 

## Features
- User authentication (session-based)

- Project submission with optional image uploads

- Secure image validation (size & MIME type)

- Admin approval workflow

- Public and admin-only views

- Flash messaging for user feedback

- PostgreSQL-backed persistence

- Clean separation of handlers, repositories, and models

## Tech Stack
- **Go** (net/http)

- **PostgreSQL**

- **pgx / pgxpool**

- **scs** (session management)

- **html/template**

- **Docker Compose** (local database)

- **UUID-based file naming**

## Architecture Overview
    /cmd

    /internal

      /handlers        HTTP handlers
      /repository      Database access layer
      /models          Domain models
      /middleware      Auth & authorization

    /uploads           User-uploaded images

    /templates         HTML templates

**The application follows a clear separation of concerns:**

- Handlers deal with HTTP

- Repositories encapsulate SQL

- Models define domain data

- Middleware handles cross-cutting concerns

## Structured Logging & Observability

The application implements structured, JSON-based logging using Go’s `slog` package.
Each incoming HTTP request is assigned a unique `request_id` and logged consistently across middleware and handlers, enabling reliable request tracing and correlation.
Logs include contextual metadata such as service name, environment, version, HTTP method, path, status code, duration, and authenticated user where applicable.
Application logs are written to disk and shipped via a lightweight log forwarder into OpenSearch, where they can be explored and visualized using OpenSearch Dashboards. 

This design prepares the application for production-grade observability pipelines (e.g. ingestion via Filebeat into OpenSearch/Elastic) and enables efficient filtering, debugging, and monitoring at scale.

## Getting Started (Local Development)
### Prerequisites
- Go 1.22+

- Docker & Docker Compose

### 1. Clone the repository
`git clone https://github.com/janphilippgutt/ProjectHub.git`

`cd projecthub`

### 2. Environment variables
`cp .env.example .env
`
Edit `.env` with your local configuration.

### 3. Start PostgeSQL
`docker compose up -d`

### 4. Run migrations
`psql -h localhost -p 5433 -U your_user -d projecthub < schema.sql`

### 5. Run the application
`go run main.go`

Visit: http://localhost:8080 (or the port you specified in .env respectively)

## Security Considerations
- Environment variables are used for all secrets

- Uploaded files are:

    - Size-limited

    - MIME-type validated

    - Stored outside templates

- Session data is server-side

- Admin routes are protected by authorization middleware

## Status

This project is under active development and serves as a production-style backend portfolio project.

## Licence
MIT