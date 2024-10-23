# Project Structure and Setup Guide

```
contractor-management/
├── .env.example                # Environment variables template
├── Dockerfile                  # Main application Dockerfile
├── docker-compose.yml         # Docker compose for local development
├── Makefile                   # Build and development commands
├── README.md                  # Project documentation
├── cmd/
│   └── api/
│       └── main.go            # Application entry point
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── domain/
│   │   ├── contractor/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── location/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── visit/
│   │       ├── entity.go
│   │       ├── repository.go
│   │       └── service.go
│   ├── infrastructure/
│   │   ├── persistence/
│   │   │   ├── postgres/
│   │   │   └── redis/
│   │   └── s3/
│   ├── middleware/
│   │   ├── auth.go
│   │   └── ratelimit.go
│   └── workflow/
│       ├── activities/
│       └── workflows/
├── migrations/
│   └── postgres/
└── pkg/
    ├── logger/
    └── validator/
```