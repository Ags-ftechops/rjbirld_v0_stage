# Contractor Management System

A secure and scalable backend system for managing contractors, locations, and document workflows.

## Prerequisites

- Docker and Docker Compose
- Go 1.21 or later (for local development)
- Make

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/your-org/contractor-management
cd contractor-management
```

2. Copy the environment file and configure it:
```bash
cp .env.example .env
# Edit .env with your configurations
```

3. Start the services:
```bash
make docker-up
```

4. Run database migrations:
```bash
make migrate-up
```

The application will be available at:
- API: http://localhost:8080
- Temporal Web UI: http://localhost:8088

## Project Structure

- `cmd/api`: Application entry point
- `config`: Configuration management
- `internal`: Internal application code
  - `domain`: Business logic and entities
  - `infrastructure`: External services integration
  - `middleware`: HTTP middleware
  - `workflow`: Temporal workflows and activities
- `migrations`: Database migration files
- `pkg`: Shared packages

## Development

1. Run tests:
```bash
make test
```

2. Build the application:
```bash
make build
```

3. Clean up:
```bash
make clean
```

## API Documentation

The API documentation is available at `/swagger/index.html` when running in development mode.

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request